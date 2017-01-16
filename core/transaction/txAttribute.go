package transaction

import (
	"GoOnchain/common/serialization"
	"io"
)

//TxAttribute descirbte the specific attributes of transcation
type TxAttribute struct {
	Usage TransactionAttributeUsage
	Date  []byte
	Size  uint32
}

func NewTxAttribute(u TransactionAttributeUsage, d []byte) TxAttribute {
	tx := TxAttribute{u, d, 0}
	tx.Size = uint32(tx.GetSize())
	return tx
}

func (u *TxAttribute) GetSize() uint64 {
	//TODO: DescriptionUrl only,other TransactionAttributeUsage type need to be impelment
	if u.Usage == DescriptionUrl {
		return uint64(len([]byte{(byte(0xff))}) + len([]byte{(byte(0xff))}) + len(u.Date))
	}
	return 0
}

func (tx *TxAttribute) Serialize(w io.Writer) {
	//TODO: DescriptionUrl only,other TransactionAttributeUsage type need to be impelment
	serialization.WriteVarBytes(w, []byte{byte(tx.Usage)})
	if tx.Usage == DescriptionUrl {
		serialization.WriteVarBytes(w, tx.Date)
	}
}

func (tx *TxAttribute) Deserialize(r io.Reader) error {
	//TODO: DescriptionUrl only,other TransactionAttributeUsage type need to be impelment
	val, _ := serialization.ReadVarBytes(r)
	tx.Usage = TransactionAttributeUsage(val[0])
	if tx.Usage == DescriptionUrl {
		tx.Date, _ = serialization.ReadVarBytes(r)
	}
	return nil
}

type TransactionAttributeUsage byte

const (
	ContractHash   TransactionAttributeUsage = 0x00
	ECDH02         TransactionAttributeUsage = 0x02 //用于ECDH密钥交换的公钥，该公钥的第一个字节为0x02
	ECDH03         TransactionAttributeUsage = 0x03 //用于ECDH密钥交换的公钥，该公钥的第一个字节为0x03
	Script         TransactionAttributeUsage = 0x20 //用于对交易进行额外的验证
	Vote           TransactionAttributeUsage = 0x30
	DescriptionUrl TransactionAttributeUsage = 0x81
	Description    TransactionAttributeUsage = 0x90

	Hash1  TransactionAttributeUsage = 0xa1
	Hash2  TransactionAttributeUsage = 0xa2
	Hash3  TransactionAttributeUsage = 0xa3
	Hash4  TransactionAttributeUsage = 0xa4
	Hash5  TransactionAttributeUsage = 0xa5
	Hash6  TransactionAttributeUsage = 0xa6
	Hash7  TransactionAttributeUsage = 0xa7
	Hash8  TransactionAttributeUsage = 0xa8
	Hash9  TransactionAttributeUsage = 0xa9
	Hash10 TransactionAttributeUsage = 0xaa
	Hash11 TransactionAttributeUsage = 0xab
	Hash12 TransactionAttributeUsage = 0xac
	Hash13 TransactionAttributeUsage = 0xad
	Hash14 TransactionAttributeUsage = 0xae
	Hash15 TransactionAttributeUsage = 0xaf
	/// 备注
	Remark   TransactionAttributeUsage = 0xf0
	Remark1  TransactionAttributeUsage = 0xf1
	Remark2  TransactionAttributeUsage = 0xf2
	Remark3  TransactionAttributeUsage = 0xf3
	Remark4  TransactionAttributeUsage = 0xf4
	Remark5  TransactionAttributeUsage = 0xf5
	Remark6  TransactionAttributeUsage = 0xf6
	Remark7  TransactionAttributeUsage = 0xf7
	Remark8  TransactionAttributeUsage = 0xf8
	Remark9  TransactionAttributeUsage = 0xf9
	Remark10 TransactionAttributeUsage = 0xfa
	Remark11 TransactionAttributeUsage = 0xfb
	Remark12 TransactionAttributeUsage = 0xfc
	Remark13 TransactionAttributeUsage = 0xfd
	Remark14 TransactionAttributeUsage = 0xfe
	Remark15 TransactionAttributeUsage = 0xff
)
