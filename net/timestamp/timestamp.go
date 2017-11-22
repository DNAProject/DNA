package timestamp

import (
	"DNA/common/log"
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"github.com/phayes/cryptoid"
	"io/ioutil"
	"math/big"
	"mime"
	"net/http"
	"time"

	_ "crypto/sha1"   // Register hash function
	_ "crypto/sha256" // Register hash function
	_ "crypto/sha512" // Register hash function
)

func init() {
	err1 := mime.AddExtensionType(".tsq", "application/timestamp-query")
	err2 := mime.AddExtensionType(".tsr", "application/timestamp-reply")
	if err1 != nil || err2 != nil {
		panic(errors.New("rfc3161: failed to register mime type"))
	}
}

// Misc Errors
var (
	ErrUnrecognizedData = errors.New("rfc3161: Got unrecognized data and end of DER")
)

// OID Identifiers
var (
	// RFC-5280: { id-kp 8 }
	// RFC-3161: {iso(1) identified-organization(3) dod(6) internet(1) security(5) mechanisms(5) pkix(7) kp (3) timestamping (8)}
	OidExtKeyUsageTimeStamping = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 8}

	// Certificate extension: "extKeyUsage": {joint-iso-itu-t(2) ds(5) certificateExtension(29) extKeyUsage(37)}
	OidExtKeyUsage = asn1.ObjectIdentifier{2, 5, 29, 37}

	// RFC-5652: Content Type: {iso(1) member-body(2) us(840) rsadsi(113549) pkcs(1) pkcs-9(9) contentType(3)}
	OidContentType = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 3}

	// RFC-5652: Message Digest: {iso(1) member-body(2) us(840) rsadsi(113549) pkcs(1) pkcs-9(9) messageDigest(4)}
	OidMessageDigest = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 4}

	// RFC-5652: iso(1) member-body(2) us(840) rsadsi(113549) pkcs(1) pkcs7(7) 2
	OidSignedData = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 2}

	// RFC-3161: iso(1) member-body(2) us(840) rsadsi(113549) pkcs(1) pkcs-9(9) smime(16) ct(1) 4
	OidContentTypeTSTInfo = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 16, 1, 4}
)

// Supported Extensions.
var supportedExtensions []asn1.ObjectIdentifier

// RootCerts is any additional trusted root certificates.
// It should only be used for testing.
// It must be initialized with x509.NewCertPool
var RootCerts *x509.CertPool

// RegisterExtension registers a supported Extension.
// This is intended to be called from the init function in
// packages that implement support for these extensions.
// A TimeStampReq or TimeStampResp with an unregistered
// critical extension will return an error when verified.
func RegisterExtension(extension asn1.ObjectIdentifier) {
	if supportedExtensions == nil {
		supportedExtensions = make([]asn1.ObjectIdentifier, 0, 0)
	}

	// Check if it already exists
	for _, ext := range supportedExtensions {
		if ext.Equal(extension) {
			return
		}
	}

	// Add it
	supportedExtensions = append(supportedExtensions, extension)
}

// ListExtensions lists all supported extensions
func ListExtensions() []asn1.ObjectIdentifier {
	if supportedExtensions == nil {
		return make([]asn1.ObjectIdentifier, 0, 0)
	} else {
		return supportedExtensions
	}
}

// PKIStatusInfo contains complete information about the status of the Time Stamp Response
type PKIStatusInfo struct {
	Status       PKIStatus
	StatusString string         `asn1:"optional,utf8"`
	FailInfo     PKIFailureInfo `asn1:"optional"`
}

func (si *PKIStatusInfo) Error() string {
	var output string
	output += si.Status.Error()
	if si.Status.IsError() {
		output += ": " + si.FailInfo.Error()
	}
	if si.StatusString != "" {
		output += ": " + si.StatusString
	}
	return output
}

// PKIStatus carries the specific status code about the status of the Time Stamp Response.
type PKIStatus int

// When the status contains the value zero or one, a TimeStampToken MUST
// be present.  When status contains a value other than zero or one, a
// TimeStampToken MUST NOT be present.  One of the following values MUST
//  be contained in status
const (
	StatusGranted                = iota // When the PKIStatus contains the value zero a TimeStampToken, as requested, is present.
	StatusGrantedWithMods               // When the PKIStatus contains the value one a TimeStampToken, with modifications, is present.
	StatusRejection                     // When the request is invalid or otherwise rejected.
	StatusWaiting                       // When the request is being processed and the client should check back later.
	StatusRevocationWarning             // Warning that a revocation is imminent.
	StatusRevocationNotification        // Notification that a revocation has occurred.
)

// IsError checks if the given Status is an error
func (status PKIStatus) IsError() bool {
	return (status != StatusGranted && status != StatusGrantedWithMods)
}

func (status PKIStatus) Error() string {
	switch status {
	case StatusGranted:
		return "A TimeStampToken, as requested, is present"
	case StatusGrantedWithMods:
		return "A TimeStampToken, with modifications, is present"
	case StatusRejection:
		return "The request is invalid or otherwise rejected"
	case StatusWaiting:
		return "The request is being processed and the client should check back later"
	case StatusRevocationWarning:
		return "A revocation is imminent"
	case StatusRevocationNotification:
		return "A revocation has occurred"
	default:
		return "Invalid PKIStatus"
	}
}

// PKIFailureInfo as defined by RFC 3161 2.4.2
type PKIFailureInfo int

// When the TimeStampToken is not present, the failInfo indicates the reason why the time-stamp
// request was rejected and may be one of the following values.
const (
	FailureBadAlg               = 0  // Unrecognized or unsupported Algorithm Identifier.
	FailureBadRequest           = 2  // Transaction not permitted or supported.
	FailureDataFormat           = 5  // The data submitted has the wrong format.
	FailureTimeNotAvailabe      = 14 // The TSA's time source is not available.
	FailureUnacceptedPolicy     = 15 // The requested TSA policy is not supported by the TSA.
	FailureUunacceptedExtension = 16 // The requested extension is not supported by the TSA.
	FailureAddInfoNotAvailable  = 17 // The additional information requested could not be understood or is not available.
	FailureSystemFailure        = 25 // The request cannot be handled due to system failure.
)

func (fi PKIFailureInfo) Error() string {
	switch fi {
	case FailureBadAlg:
		return "Unrecognized or unsupported Algorithm Identifier"
	case FailureBadRequest:
		return "Transaction not permitted or supported"
	case FailureDataFormat:
		return "The data submitted has the wrong format"
	case FailureTimeNotAvailabe:
		return "The TSA's time source is not available"
	case FailureUnacceptedPolicy:
		return "The requested TSA policy is not supported by the TSA"
	case FailureUunacceptedExtension:
		return "The requested extension is not supported by the TSA"
	case FailureAddInfoNotAvailable:
		return "The additional information requested could not be understood or is not available"
	case FailureSystemFailure:
		return "The request cannot be handled due to system failure"
	default:
		return "Invalid PKIFailureInfo"
	}
}

// Errors
var (
	ErrInvalidDigestSize = errors.New("rfc3161: Invalid Message Digest. Invalid size for the given hash algorithm")
	ErrUnsupportedHash   = errors.New("rfc3161: Unsupported Hash Algorithm")
	ErrUnsupportedExt    = errors.New("rfc3161: Unsupported Critical Extension")
)

// TimeStampReq contains a full Time Stamp Request as defined by RFC 3161
// It is also known as a "Time Stamp Query"
// When stored into a file it should contain the extension ".tsq"
// It has a mime-type of "application/timestamp-query"
type TimeStampReq struct {
	Version        int                   `asn1:"default:1"`
	MessageImprint MessageImprint        // A hash algorithm OID and the hash value of the data to be time-stamped
	ReqPolicy      asn1.ObjectIdentifier `asn1:"optional"` // Identifier for the policy. For many TSA's, often the same as SignedData.DigestAlgorithm
	Nonce          *big.Int              `asn1:"optional"` // Nonce could be up to 160 bits
	CertReq        bool                  `asn1:"optional"` // If set to true, the TSA's certificate MUST be provided in the response.
	Extensions     []pkix.Extension      `asn1:"optional,tag:0"`
}

// MessageImprint contains hash algorithm OID and the hash digest of the data to be time-stamped
type MessageImprint struct {
	HashAlgorithm pkix.AlgorithmIdentifier
	HashedMessage []byte
}

// NewTimeStampReq creates a new Time Stamp Request, given a crypto.Hash algorithm and a message digest
func NewTimeStampReq(hash crypto.Hash, digest []byte) (*TimeStampReq, error) {
	tsr := new(TimeStampReq)
	tsr.Version = 1

	err := tsr.SetHashDigest(hash, digest)
	if err != nil {
		return nil, err
	}

	return tsr, nil
}

// ReadTSQ reads a .tsq file into a TimeStampReq
func ReadTSQ(filename string) (*TimeStampReq, error) {
	der, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	req := new(TimeStampReq)
	rest, err := asn1.Unmarshal(der, req)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return req, ErrUnrecognizedData
	}
	return req, nil
}

// SetHashDigest sets the Hash Algorithm and the Hash Digest for the Time Stamp Request
func (tsr *TimeStampReq) SetHashDigest(hash crypto.Hash, digest []byte) error {
	if len(digest) != hash.Size() {
		return ErrInvalidDigestSize
	}
	pkixAlgo := pkix.AlgorithmIdentifier{
		Algorithm: cryptoid.HashAlgorithmByCrypto(hash).OID,
	}

	tsr.MessageImprint.HashAlgorithm = pkixAlgo
	tsr.MessageImprint.HashedMessage = digest

	return nil
}

// GetHash will get the crypto.Hash for the Time Stamp Request
// The Hash will be 0 if it is not recognized
func (tsr *TimeStampReq) GetHash() crypto.Hash {
	hashAlgo, err := cryptoid.HashAlgorithmByOID(tsr.MessageImprint.HashAlgorithm.Algorithm.String())
	if err != nil {
		return 0
	}
	return hashAlgo.Hash
}

// GenerateNonce generates a 128 bit nonce for the Time Stamp Request
// If a different size is required then set manually with tsr.Nonce.SetBytes()
func (tsr *TimeStampReq) GenerateNonce() error {
	// Generate a 128 bit nonce
	b := make([]byte, 16, 16)

	_, err := rand.Read(b)
	if err != nil {
		return err
	}

	tsr.Nonce = new(big.Int)
	tsr.Nonce.SetBytes(b)

	return nil
}

// Verify does a basic sanity check of the Time Stamp Request
// Checks to make sure the hash is supported, the digest matches the hash,
// and no unsupported critical extensions exist. Be sure to add all supported
// extentions to rfc3161.SupportedExtensions.
func (tsr *TimeStampReq) Verify() error {
	hash := tsr.GetHash()
	if hash == 0 {
		return ErrUnsupportedHash
	}
	if len(tsr.MessageImprint.HashedMessage) != hash.Size() {
		return ErrInvalidDigestSize
	}

	// Check for any unsupported critical extensions
	// Critical Extensions should be registered in rfc3161.SupportedExtensions
	if tsr.Extensions != nil {
		for _, ext := range tsr.Extensions {
			if ext.Critical {
				supported := false
				if supportedExtensions != nil {
					for _, se := range supportedExtensions {
						if se.Equal(ext.Id) {
							supported = true
							break
						}
					}
				}
				if !supported {
					return ErrUnsupportedExt
				}
			}
		}
	}

	return nil
}

// Errors
var (
	ErrIncorrectNonce              = errors.New("rfc3161: response: Response has incorrect nonce")
	ErrNoTST                       = errors.New("rfc3161: response: Response does not contain TSTInfo")
	ErrNoCertificate               = errors.New("rfc3161: response: No certificates provided")
	ErrMismatchedCertificates      = errors.New("rfc3161: response: Mismatched certificates")
	ErrCertificateKeyUsage         = errors.New("rfc3161: response: certificate: Invalid KeyUsage field")
	ErrCertificateExtKeyUsageUsage = errors.New("rfc3161: response: certificate: Invalid ExtKeyUsage field")
	ErrCertificateExtension        = errors.New("rfc3161: response: certificate: Missing critical timestamping extension")
	ErrInvalidSignatureDigestAlgo  = errors.New("rfc3161: response: Invalid signature digest algorithm")
	ErrUnsupportedSignerInfos      = errors.New("rfc3161: response: package only supports responses with a single SignerInfo")
	ErrUnableToParseSID            = errors.New("rfc3161: response: Unable to parse SignerInfo.sid")
	ErrVerificationError           = errors.New("rfc3161: response: Verfication error")
	ErrInvalidOID                  = errors.New("rfc3161: response: Invalid OID")
)

// TimeStampResp contains a full Time Stamp Response as defined by RFC 3161
// It is also known as a "Time Stamp Reply"
// When stored into a file it should contain the extension ".tsr"
// It has a mime-type of "application/timestamp-reply"
type TimeStampResp struct {
	Status         PKIStatusInfo
	TimeStampToken `asn1:"optional"`
}

// ReadTSR reads a .tsr file into a TimeStampResp
func ReadTSR(filename string) (*TimeStampResp, error) {
	der, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	resp := new(TimeStampResp)
	rest, err := asn1.Unmarshal(der, resp)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return resp, ErrUnrecognizedData
	}
	return resp, nil
}

func (resp *TimeStampResp) Verify(req *TimeStampReq, cert *x509.Certificate) error {
	tst, err := resp.GetTSTInfo()
	if err != nil {
		return err
	}

	// Verify the request for sanity's sake
	err = req.Verify()
	if err != nil {
		return err
	}

	// Verify the status
	if resp.Status.Status.IsError() {
		return &resp.Status
	}

	// Verify the nonce
	if req.Nonce == nil || tst.Nonce == nil {
		if req.Nonce != tst.Nonce {
			return ErrIncorrectNonce
		}
	} else if req.Nonce.Cmp(tst.Nonce) != 0 {
		return ErrIncorrectNonce
	}

	// Verify that the OIDs are correct
	if !resp.ContentType.Equal(OidSignedData) || !resp.EContentType.Equal(OidContentTypeTSTInfo) {
		return ErrInvalidOID
	}

	// Get the certificate
	respcert, err := resp.GetSigningCert()
	if err != nil {
		return err
	}
	// Rationalize the passed-in certificate vis-a-vis certificate in the response
	if req.CertReq {
		if respcert != nil && cert != nil {
			if !bytes.Equal(cert.Raw, respcert.Raw) {
				return ErrMismatchedCertificates
			}
		} else if cert == nil {
			cert = respcert
		}
	}
	if cert == nil {
		return ErrNoCertificate
	}

	// Get any intermediates that might be needed
	intermediates, err := resp.GetCertificates()
	if err != nil && err != ErrNoCertificate {
		return err
	}
	interpool := x509.NewCertPool()
	for _, intercert := range intermediates {
		interpool.AddCert(intercert)
	}

	// Verify the certificate
	err = resp.VerifyCertificate(cert, interpool)
	if err != nil {
		return err
	}

	// Verify the signature
	err = resp.VerifySignature(cert)
	if err != nil {
		return err
	}

	return nil
}

func (resp *TimeStampResp) VerifyCertificate(cert *x509.Certificate, intermediates *x509.CertPool) error {
	if cert == nil {
		return ErrNoCertificate
	}

	// Key usage must contain the KeyUsageDigitalSignature bit
	// and MAY contain the non-repudiation / content-commitment bit
	if cert.KeyUsage != x509.KeyUsageDigitalSignature && cert.KeyUsage != (x509.KeyUsageDigitalSignature+x509.KeyUsageContentCommitment) {
		return ErrCertificateKeyUsage
	}

	// Next check the extended key usage
	// Only one ExtKeyUsage may be defined as per RFC 3161
	if len(cert.ExtKeyUsage) != 1 {
		return ErrCertificateExtKeyUsageUsage
	}
	if cert.ExtKeyUsage[0] != x509.ExtKeyUsageTimeStamping {
		return ErrCertificateExtKeyUsageUsage
	}

	// Check to make sure it has the correct extension
	// Only one Extended Key Usage may be defined, it must be critical,
	// and it must be OidExtKeyUsageTimeStamping
	for _, ext := range cert.Extensions {
		if ext.Id.Equal(OidExtKeyUsage) {
			if !ext.Critical {
				return ErrCertificateExtKeyUsageUsage
			}
			var rfc3161Ext []asn1.ObjectIdentifier
			_, err := asn1.Unmarshal(ext.Value, &rfc3161Ext)
			if err != nil {
				return err
			}
			if len(rfc3161Ext) != 1 {
				return ErrCertificateExtKeyUsageUsage
			}
			if !rfc3161Ext[0].Equal(OidExtKeyUsageTimeStamping) {
				return ErrCertificateExtKeyUsageUsage
			}
		}
	}

	// Verify the certificate chain
	opts := x509.VerifyOptions{
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageTimeStamping},
		Roots:         RootCerts,
		Intermediates: intermediates,
	}
	_, err := cert.Verify(opts)
	if err != nil {
		return err
	}

	return nil
}

// TimeStampToken is a wrapper than contains the OID for a TimeStampToken
// as well as the wrapped SignedData
type TimeStampToken struct {
	ContentType asn1.ObjectIdentifier // MUST BE OidSignedData
	SignedData  `asn1:"tag:0,explicit,optional"`
}

// SignedData is a shared-standard as defined by RFC 2630
type SignedData struct {
	Version          int                        `asn1:"default:4"`
	DigestAlgorithms []pkix.AlgorithmIdentifier `asn1:"set"`
	EncapsulatedContentInfo
	Certificates asn1.RawValue          `asn1:"optional,set,tag:0"` // Certificate DER. Use GetCertificates() to get the x509.Certificate list
	CRLs         []pkix.CertificateList `asn1:"optional,tag:1"`
	SignerInfos  []SignerInfo           `asn1:"set"`
}

// GetSigningCert gets the signer and the associated certificate
// The certificate may be nil if the request did not ask for it
func (sd *SignedData) GetSigningCert() (*x509.Certificate, error) {
	// Get the signerInfo
	if len(sd.SignerInfos) != 1 {
		return nil, ErrUnsupportedSignerInfos
	}
	signer := sd.SignerInfos[0]
	id, err := signer.GetSID()
	if err != nil {
		return nil, err
	}

	var cert *x509.Certificate
	if len(sd.Certificates.Bytes) != 0 {
		certs, err := x509.ParseCertificates(sd.Certificates.Bytes)
		if err != nil {
			return nil, err
		}
		for _, checkcert := range certs {
			switch sid := id.(type) {
			case *IssuerAndSerialNumber:
				if checkcert.SerialNumber.Cmp(sid.SerialNumber) == 0 {
					cert = checkcert
					break
				}
			case []byte:
				if bytes.Equal(checkcert.SubjectKeyId, sid) {
					cert = checkcert
					break
				}
			default:
				return nil, ErrUnableToParseSID
			}
		}
	}
	return cert, nil
}

// VerifySignature verifies that the given certificate signed the TSTInfo
func (sd *SignedData) VerifySignature(cert *x509.Certificate) error {
	// Get the signerInfo
	if len(sd.SignerInfos) != 1 {
		return ErrUnsupportedSignerInfos
	}
	signer := sd.SignerInfos[0]

	hashAlgo, err := cryptoid.HashAlgorithmByOID(signer.DigestAlgorithm.Algorithm.String())
	if err != nil {
		return err
	}

	// Marshal the Signed Attributes
	derbytes, err := asn1.Marshal(signer.SignedAttrs)
	if err != nil {
		return err
	}

	// Hack the DER bytes of the Signed Attributes to be EXPLICIT SET
	derbytes[0] = 0x31
	derbytes[1] = 0x81

	// Hash the DER bytes
	hash := hashAlgo.Hash.New()
	hash.Write(derbytes)
	digest := hash.Sum(nil)

	// Unpack the public key
	pub := cert.PublicKey.(*rsa.PublicKey)

	// Verify the signature
	err = rsa.VerifyPKCS1v15(pub, hashAlgo.Hash, digest, signer.Signature)
	if err != nil {
		return ErrVerificationError
	}

	// Verify the signed attributes
	// This will check the following:
	// - the content-type is of the type TSTInfo
	// - the message-digest corresponds to TSTInfo
	// - there is exactly one digest attribute and one content-type attribute
	var digestOK, contentOK bool
	var count int
	for _, attr := range signer.SignedAttrs {
		if attr.Type.Equal(OidContentType) {
			count++
			oiddata, _ := asn1.Marshal(OidContentTypeTSTInfo)
			if bytes.Equal(oiddata, attr.Value.Bytes) {
				contentOK = true
			}
		}
		if attr.Type.Equal(OidMessageDigest) {
			count++
			hash := hashAlgo.Hash.New()
			hash.Write(sd.EContent)
			digest := hash.Sum(nil)
			if bytes.Equal(digest, attr.Value.Bytes[2:]) {
				digestOK = true
			}
		}
	}
	if !digestOK || !contentOK || count != 2 {
		return ErrVerificationError
	}

	// Everything is OK
	return nil
}

// GetCertificates gets a list of x509.Certificate objects from the DER encoded Certificates field
func (sd *SignedData) GetCertificates() ([]*x509.Certificate, error) {
	if len(sd.Certificates.Bytes) == 0 {
		return nil, ErrNoCertificate
	}
	return x509.ParseCertificates(sd.Certificates.Bytes)
}

// SignerInfo is a shared-standard as defined by RFC 2630
type SignerInfo struct {
	Version            int           `asn1:"default:1"`
	SID                asn1.RawValue // CHOICE. See SignerInfo.GetSID()
	DigestAlgorithm    pkix.AlgorithmIdentifier
	SignedAttrs        []Attribute `asn1:"tag:0"`
	SignatureAlgorithm pkix.AlgorithmIdentifier
	Signature          []byte
	UnsignedAtrributes []Attribute `asn1:"optional,tag:1"`
}

// GetSID Gets the certificate identifier
// It returns an interface that could be one of:
//  - *rfc3161.IssuerAndSerialNumber
//  - []byte if the identifier is a SubjectKeyId
func (sd *SignerInfo) GetSID() (interface{}, error) {
	var sid interface{}
	switch sd.Version {
	case 1:
		sid = &IssuerAndSerialNumber{}
	case 3:
		sid = []byte{}
	default:
		return nil, errors.New("Invalid SignerInfo.SID")
	}

	_, err := asn1.Unmarshal(sd.SID.FullBytes, sid)
	if err != nil {
		return nil, err
	}
	return sid, nil
}

// IssuerAndSerialNumber is defined in RFC 2630
type IssuerAndSerialNumber struct {
	IssuerName   pkix.RDNSequence
	SerialNumber *big.Int
}

// Attribute is defined in RFC 2630
// The fields of type SignedAttribute and UnsignedAttribute have the
// following meanings:
//
//   Type indicates the type of attribute.  It is an object
//   identifier.
//
//   Value is a set of values that comprise the attribute.  The
//   type of each value in the set can be determined uniquely by
//   Type.
type Attribute struct {
	Type  asn1.ObjectIdentifier
	Value asn1.RawValue
}

// EncapsulatedContentInfo is defined in RFC 2630
//
// The fields of type EncapsulatedContentInfo of the SignedData
// construct have the following meanings:
//
// eContentType is an object identifier that uniquely specifies the
// content type.  For a time-stamp token it is defined as:
//
// id-ct-TSTInfo  OBJECT IDENTIFIER ::= { iso(1) member-body(2)
// us(840) rsadsi(113549) pkcs(1) pkcs-9(9) smime(16) ct(1) 4}
//
// eContent is the content itself, carried as an octet string.
// The eContent SHALL be the DER-encoded value of TSTInfo.
//
// The time-stamp token MUST NOT contain any signatures other than the
// signature of the TSA.  The certificate identifier (ESSCertID) of the
// TSA certificate MUST be included as a signerInfo attribute inside a
// SigningCertificate attribute.
type EncapsulatedContentInfo struct {
	EContentType asn1.ObjectIdentifier // MUST BE OidContentTypeTSTInfo
	EContent     asn1.RawContent       `asn1:"explicit,optional,tag:0"` // DER encoding of TSTInfo
}

// GetTSTInfo unpacks the DER encoded TSTInfo
func (eci *EncapsulatedContentInfo) GetTSTInfo() (*TSTInfo, error) {
	if len(eci.EContent) == 0 {
		return nil, ErrNoTST
	}

	tstinfo := new(TSTInfo)
	rest, err := asn1.Unmarshal(eci.EContent, tstinfo)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return tstinfo, ErrUnrecognizedData
	}

	return tstinfo, nil
}

// TSTInfo is the acutal DER signed data and represents the core of the Time Stamp Reponse.
// It contains the time-stamp, the accuracy, and all other pertinent informatuon
type TSTInfo struct {
	Version        int                   `json:"version" asn1:"default:1"`
	Policy         asn1.ObjectIdentifier `json:"policy"`                           // Identifier for the policy. For many TSA's, often the same as SignedData.DigestAlgorithm
	MessageImprint MessageImprint        `json:"message-imprint"`                  // MUST have the same value of MessageImprint in matching TimeStampReq
	SerialNumber   *big.Int              `json:"serial-number"`                    // Time-Stamping users MUST be ready to accommodate integers up to 160 bits
	GenTime        time.Time             `json:"gen-time"`                         // The time at which it was stamped
	Accuracy       Accuracy              `json:"accuracy" asn1:"optional"`         // Accuracy represents the time deviation around the UTC time.
	Ordering       bool                  `json:"ordering" asn1:"optional"`         // True if SerialNumber increases monotonically with time.
	Nonce          *big.Int              `json:"nonce" asn1:"optional"`            // MUST be present if the similar field was present in TimeStampReq.  In that case it MUST have the same value.
	TSA            asn1.RawValue         `json:"tsa" asn1:"optional,tag:0"`        // This is a CHOICE (See RFC 3280 for all choices). See https://github.com/golang/go/issues/13999 for information on handling.
	Extensions     []pkix.Extension      `json:"extensions" asn1:"optional,tag:1"` // List of extensions
}

// Accuracy represents the time deviation around the UTC time.
//
// If either seconds, millis or micros is missing, then a value of zero
// MUST be taken for the missing field.
//
// By adding the accuracy value to the GeneralizedTime, an upper limit
// of the time at which the time-stamp token has been created by the TSA
// can be obtained.  In the same way, by subtracting the accuracy to the
// GeneralizedTime, a lower limit of the time at which the time-stamp
// token has been created by the TSA can be obtained.
//
// Accuracy can be decomposed in seconds, milliseconds (between 1-999)
// and microseconds (1-999), all expressed as integer.
//
// When the accuracy field is not present, then the accuracy
// may be available through other means, e.g., the TSAPolicyId.
type Accuracy struct {
	Seconds int `asn1:"optional"`
	Millis  int `asn1:"optional,tag:0"`
	Micros  int `asn1:"optional,tag:1"`
}

// Duration gets the time.Duration representation of the Accuracy
func (acc *Accuracy) Duration() time.Duration {
	return (time.Duration(acc.Seconds) * time.Second) + (time.Duration(acc.Millis) * time.Millisecond) + (time.Duration(acc.Micros) * time.Microsecond)
}

type TimeStampClient struct {
	HTTPClient *http.Client
	URL        string
}

//digest: sha256(message)
func (self *TimeStampClient) FetchTimeStampToken(digest []byte) ([]byte, int64, error) {
	req, err := NewTimeStampReq(crypto.SHA256, digest)
	if err != nil {
		return nil, 0, err
	}
	req.CertReq = true

	err = req.GenerateNonce()
	if err != nil {
		return nil, 0, err
	}

	err = req.Verify()
	if err != nil {
		return nil, 0, err
	}

	resp, err := self.Do(req)
	if err != nil {
		log.Error("timestampclient: cannot get response:", err)
		return nil, 0, err
	}

	err = resp.Verify(req, nil)
	if err != nil {
		log.Error("timestampclient: verify error:", err)
		return nil, 0, err
	}

	data, err := asn1.Marshal(resp.TimeStampToken)
	if err != nil {
		log.Error("timestampclient: failed to marshal timestamptoken:", err)
		return nil, 0, err
	}

	info, _ := resp.GetTSTInfo()

	return data, info.GenTime.Unix(), nil
}

// NewClient creates a new rfc3161.Client given a URL.
func NewClient(url string) *TimeStampClient {
	client := new(TimeStampClient)
	client.HTTPClient = http.DefaultClient
	client.URL = url
	return client
}

// Do a time stamp request and get back the Time Stamp Response.
// This will not verify the response. It is the caller's responsibility
// to call resp.Verify() on the returned TimeStampResp.
func (client *TimeStampClient) Do(tsq *TimeStampReq) (*TimeStampResp, error) {
	der, err := asn1.Marshal(*tsq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", client.URL, bytes.NewBuffer(der))
	if err != nil {
		return nil, err
	}

	//disable reuse transport to avoid EOF Error
	// see: https://stackoverflow.com/questions/17714494/golang-http-request-results-in-eof-errors-when-making-multiple-requests-successi/23963271#23963271
	req.Close = true
	req.Header.Set("Content-Type", "application/timestamp-query")

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tsr := new(TimeStampResp)
	rest, err := asn1.Unmarshal(body, tsr)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, ErrUnrecognizedData
	}

	if tsr.Status.Status.IsError() {
		return tsr, &tsr.Status
	}
	return tsr, nil
}
