package common

import (
	. "DNA/common"
	. "DNA/common/config"
	"DNA/common/log"
	"DNA/core/ledger"
	tx "DNA/core/transaction"
	. "DNA/net/httpjsonrpc"
	Err "DNA/net/httprestful/error"
	. "DNA/net/protocol"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var node Noder
var pushBlockFlag bool = true
var oauthClient = NewOauthClient()

type ApiServer interface {
	Start() error
	Stop()
}

func SetNode(n Noder) {
	node = n
}
func CheckPushBlock() bool {
	return pushBlockFlag
}

//Node
func GetConnectionCount(cmd map[string]interface{}) map[string]interface{} {
	resp := ResponsePack(Err.SUCCESS)
	if node != nil {
		resp["Result"] = node.GetConnectionCnt()
	}

	return resp
}

//Block
func GetBlockHeight(cmd map[string]interface{}) map[string]interface{} {
	resp := ResponsePack(Err.SUCCESS)
	resp["Result"] = ledger.DefaultLedger.Blockchain.BlockHeight
	return resp
}
func GetBlockHash(cmd map[string]interface{}) map[string]interface{} {
	resp := ResponsePack(Err.SUCCESS)
	param := cmd["Height"].(string)
	if len(param) == 0 {
		resp["Error"] = Err.INVALID_PARAMS
		return resp
	}
	height, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		resp["Error"] = Err.INVALID_PARAMS
		return resp
	}
	hash, err := ledger.DefaultLedger.Store.GetBlockHash(uint32(height))
	if err != nil {
		resp["Error"] = Err.INVALID_PARAMS
		return resp
	}
	resp["Result"] = ToHexString(hash.ToArrayReverse())
	return resp
}
func GetBlockInfo(block *ledger.Block) BlockInfo {
	hash := block.Hash()
	blockHead := &BlockHead{
		Version:          block.Blockdata.Version,
		PrevBlockHash:    ToHexString(block.Blockdata.PrevBlockHash.ToArrayReverse()),
		TransactionsRoot: ToHexString(block.Blockdata.TransactionsRoot.ToArrayReverse()),
		Timestamp:        block.Blockdata.Timestamp,
		Height:           block.Blockdata.Height,
		ConsensusData:    block.Blockdata.ConsensusData,
		NextBookKeeper:   ToHexString(block.Blockdata.NextBookKeeper.ToArrayReverse()),
		Program: ProgramInfo{
			Code:      ToHexString(block.Blockdata.Program.Code),
			Parameter: ToHexString(block.Blockdata.Program.Parameter),
		},
		Hash: ToHexString(hash.ToArrayReverse()),
	}

	trans := make([]*Transactions, len(block.Transactions))
	for i := 0; i < len(block.Transactions); i++ {
		trans[i] = TransArryByteToHexString(block.Transactions[i])
	}

	b := BlockInfo{
		Hash:         ToHexString(hash.ToArrayReverse()),
		BlockData:    blockHead,
		Transactions: trans,
	}
	return b
}
func getBlock(hash Uint256, getTxBytes bool) (interface{}, int64) {
	block, err := ledger.DefaultLedger.Store.GetBlock(hash)
	if err != nil {
		return "", Err.UNKNOWN_BLOCK
	}
	if getTxBytes {
		w := bytes.NewBuffer(nil)
		block.Serialize(w)
		return ToHexString(w.Bytes()), Err.SUCCESS
	}
	return GetBlockInfo(block), Err.SUCCESS
}
func GetBlockByHash(cmd map[string]interface{}) map[string]interface{} {
	resp := ResponsePack(Err.SUCCESS)
	param := cmd["Hash"].(string)
	if len(param) == 0 {
		resp["Error"] = Err.INVALID_PARAMS
		return resp
	}
	var getTxBytes bool = false
	if raw, ok := cmd["Raw"].(string); ok && raw == "1" {
		getTxBytes = true
	}
	var hash Uint256
	hex, err := HexToBytesReverse(param)
	if err != nil {
		resp["Error"] = Err.INVALID_PARAMS
		return resp
	}
	if err := hash.Deserialize(bytes.NewReader(hex)); err != nil {
		resp["Error"] = Err.INVALID_TRANSACTION
		return resp
	}

	resp["Result"], resp["Error"] = getBlock(hash, getTxBytes)

	return resp
}
func GetBlockByHeight(cmd map[string]interface{}) map[string]interface{} {
	resp := ResponsePack(Err.SUCCESS)

	param := cmd["Height"].(string)
	if len(param) == 0 {
		resp["Error"] = Err.INVALID_PARAMS
		return resp
	}
	var getTxBytes bool = false
	if raw, ok := cmd["Raw"].(string); ok && raw == "1" {
		getTxBytes = true
	}
	height, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		resp["Error"] = Err.INVALID_PARAMS
		return resp
	}
	index := uint32(height)
	hash, err := ledger.DefaultLedger.Store.GetBlockHash(index)
	if err != nil {
		resp["Error"] = Err.UNKNOWN_BLOCK
		return resp
	}
	resp["Result"], resp["Error"] = getBlock(hash, getTxBytes)
	return resp
}

//Asset
func GetAssetByHash(cmd map[string]interface{}) map[string]interface{} {
	resp := ResponsePack(Err.SUCCESS)

	str := cmd["Hash"].(string)
	hex, err := HexToBytesReverse(str)
	if err != nil {
		resp["Error"] = Err.INVALID_PARAMS
		return resp
	}
	var hash Uint256
	err = hash.Deserialize(bytes.NewReader(hex))
	if err != nil {
		resp["Error"] = Err.INVALID_ASSET
		return resp
	}
	asset, err := ledger.DefaultLedger.Store.GetAsset(hash)
	if err != nil {
		resp["Error"] = Err.UNKNOWN_ASSET
		return resp
	}
	resp["Result"] = asset
	return resp
}
func GetUnspendOutput(cmd map[string]interface{}) map[string]interface{} {
	resp := ResponsePack(Err.SUCCESS)
	addr := cmd["Addr"].(string)
	assetid := cmd["Assetid"].(string)

	var programHash Uint160
	var assetHash Uint256

	bys, err := HexToBytesReverse(addr)
	if err != nil {
		resp["Error"] = Err.INVALID_PARAMS
		return resp
	}
	if err := programHash.Deserialize(bytes.NewReader(bys)); err != nil {
		resp["Error"] = Err.INVALID_PARAMS
		return resp
	}

	bys, err = HexToBytesReverse(assetid)
	if err != nil {
		resp["Error"] = Err.INVALID_PARAMS
		return resp
	}
	if err := assetHash.Deserialize(bytes.NewReader(bys)); err != nil {
		resp["Error"] = Err.INVALID_PARAMS
		return resp
	}
	type TxOutputInfo struct {
		AssetID     string
		Value       int64
		ProgramHash string
	}
	outputs := make(map[string]*TxOutputInfo)
	height := ledger.DefaultLedger.GetLocalBlockChainHeight()
	var i uint32
	// construct global UTXO table
	for i = 0; i <= height; i++ {
		block, err := ledger.DefaultLedger.GetBlockWithHeight(i)
		if err != nil {
			resp["Error"] = Err.INTERNAL_ERROR
			return resp
		}
		// skip the bookkeeping transaction
		for _, t := range block.Transactions[1:] {
			// skip the register transaction
			if t.TxType == tx.RegisterAsset {
				continue
			}
			txHash := t.Hash()
			txHashHex := ToHexString(txHash.ToArrayReverse())
			for i, output := range t.Outputs {
				if output.AssetID.CompareTo(assetHash) == 0 &&
					output.ProgramHash.CompareTo(programHash) == 0 {
					key := txHashHex + ":" + strconv.Itoa(i)
					asset := ToHexString(output.AssetID.ToArrayReverse())
					pHash := ToHexString(output.ProgramHash.ToArrayReverse())
					value := int64(output.Value)
					info := &TxOutputInfo{
						asset,
						value,
						pHash,
					}
					outputs[key] = info
				}
			}
		}
	}
	// delete spent output from global UTXO table
	height = ledger.DefaultLedger.GetLocalBlockChainHeight()
	for i = 0; i <= height; i++ {
		block, err := ledger.DefaultLedger.GetBlockWithHeight(i)
		if err != nil {
			return DnaRpcInternalError
		}
		// skip the bookkeeping transaction
		for _, t := range block.Transactions[1:] {
			// skip the register transaction
			if t.TxType == tx.RegisterAsset {
				continue
			}
			for _, input := range t.UTXOInputs {
				refer := ToHexString(input.ReferTxID.ToArrayReverse())
				index := strconv.Itoa(int(input.ReferTxOutputIndex))
				key := refer + ":" + index
				delete(outputs, key)
			}
		}
	}
	resp["Result"] = outputs
	return resp
}

//Transaction
func GetTransactionByHash(cmd map[string]interface{}) map[string]interface{} {
	resp := ResponsePack(Err.SUCCESS)

	str := cmd["Hash"].(string)
	bys, err := HexToBytesReverse(str)
	if err != nil {
		resp["Error"] = Err.INVALID_PARAMS
		return resp
	}
	var hash Uint256
	err = hash.Deserialize(bytes.NewReader(bys))
	if err != nil {
		resp["Error"] = Err.INVALID_TRANSACTION
		return resp
	}
	tx, err := ledger.DefaultLedger.Store.GetTransaction(hash)
	if err != nil {
		resp["Error"] = Err.UNKNOWN_TRANSACTION
		return resp
	}
	if raw, ok := cmd["Raw"].(string); ok && raw == "1" {
		w := bytes.NewBuffer(nil)
		tx.Serialize(w)
		resp["Result"] = ToHexString(w.Bytes())
		return resp
	}
	tran := TransArryByteToHexString(tx)
	resp["Result"] = tran
	return resp
}
func SendRawTransaction(cmd map[string]interface{}) map[string]interface{} {
	resp := ResponsePack(Err.SUCCESS)

	str, ok := cmd["Data"].(string)
	if !ok {
		resp["Error"] = Err.INVALID_PARAMS
		return resp
	}
	bys, err := HexToBytes(str)
	if err != nil {
		resp["Error"] = Err.INVALID_PARAMS
		return resp
	}
	var txn tx.Transaction
	if err := txn.Deserialize(bytes.NewReader(bys)); err != nil {
		resp["Error"] = Err.INVALID_TRANSACTION
		return resp
	}
	var hash Uint256
	hash = txn.Hash()
	if err := VerifyAndSendTx(&txn); err != nil {
		resp["Error"] = Err.INTERNAL_ERROR
		return resp
	}
	resp["Result"] = ToHexString(hash.ToArrayReverse())
	//TODO 0xd1 -> tx.InvokeCode
	if txn.TxType == 0xd1 {
		if userid, ok := cmd["Userid"].(string); ok && len(userid) > 0 {
			resp["Userid"] = userid
		}
	}
	return resp
}

//record

//config
func GetOauthServerUrl(cmd map[string]interface{}) map[string]interface{} {
	resp := ResponsePack(Err.SUCCESS)
	resp["Result"] = Parameters.OauthServerAddr
	return resp
}
func SetOauthServerUrl(cmd map[string]interface{}) map[string]interface{} {
	resp := ResponsePack(Err.SUCCESS)

	addr, ok := cmd["Addr"].(string)
	if !ok {
		resp["Error"] = Err.INVALID_PARAMS
		return resp
	}
	if len(addr) > 0 {
		var reg *regexp.Regexp
		pattern := `((http|https)://)(([a-zA-Z0-9\._-]+\.[a-zA-Z]{2,6})|([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}))(:[0-9]{1,4})*(/[a-zA-Z0-9\&%_\./-~-]*)?`
		reg = regexp.MustCompile(pattern)
		if !reg.Match([]byte(addr)) {
			resp["Error"] = Err.INVALID_PARAMS
			return resp
		}
	}
	Parameters.OauthServerAddr = addr
	resp["Result"] = Parameters.OauthServerAddr
	return resp
}
func GetNoticeServerUrl(cmd map[string]interface{}) map[string]interface{} {
	resp := ResponsePack(Err.SUCCESS)
	resp["Result"] = Parameters.NoticeServerAddr
	return resp
}

func SetPushBlockFlag(cmd map[string]interface{}) map[string]interface{} {
	resp := ResponsePack(Err.SUCCESS)
	open, ok := cmd["Open"].(bool)
	if !ok {
		resp["Error"] = Err.INVALID_PARAMS
		return resp
	}
	pushBlockFlag = open
	resp["Result"] = pushBlockFlag
	return resp
}
func SetNoticeServerUrl(cmd map[string]interface{}) map[string]interface{} {
	resp := ResponsePack(Err.SUCCESS)

	addr, ok := cmd["Addr"].(string)
	if !ok || len(addr) == 0 {
		resp["Error"] = Err.INVALID_PARAMS
		return resp
	}
	var reg *regexp.Regexp
	pattern := `((http|https)://)(([a-zA-Z0-9\._-]+\.[a-zA-Z]{2,6})|([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}))(:[0-9]{1,4})*(/[a-zA-Z0-9\&%_\./-~-]*)?`
	reg = regexp.MustCompile(pattern)
	if !reg.Match([]byte(addr)) {
		resp["Error"] = Err.INVALID_PARAMS
		return resp
	}
	Parameters.NoticeServerAddr = addr
	resp["Result"] = Parameters.NoticeServerAddr
	return resp
}
func ResponsePack(errCode int64) map[string]interface{} {
	resp := map[string]interface{}{
		"Action":  "",
		"Result":  "",
		"Error":   errCode,
		"Desc":    "",
		"Version": "1.0.0",
	}
	return resp
}

func PostRequest(cmd map[string]interface{}, url string) (map[string]interface{}, error) {

	var repMsg = make(map[string]interface{})

	data, err := json.Marshal(cmd)
	if err != nil {
		return repMsg, err
	}
	reqData := bytes.NewBuffer(data)
	transport := http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(netw, addr, time.Second*10)
			if err != nil {
				return nil, err
			}
			conn.SetDeadline(time.Now().Add(time.Second * 10))
			return conn, nil
		},
		DisableKeepAlives: false,
	}
	client := &http.Client{Transport: &transport}
	request, err := http.NewRequest("POST", url, reqData)
	if err != nil {
		return repMsg, err
	}
	request.Header.Set("Content-type", "application/json")

	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
		if response.StatusCode == 200 {
			body, _ := ioutil.ReadAll(response.Body)
			if err := json.Unmarshal(body, &repMsg); err == nil {
				return repMsg, err
			}
		}
	}

	if err != nil {
		return repMsg, err
	}

	return repMsg, err
}

func NewOauthClient() *http.Client {
	c := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*10)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 10))
				return conn, nil
			},
			DisableKeepAlives: false,
		},
	}
	return c
}

func OauthRequest(method string, cmd map[string]interface{}, url string) (map[string]interface{}, error) {

	var repMsg = make(map[string]interface{})
	var response *http.Response
	var err error
	switch method {
	case "GET":

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return repMsg, err
		}
		response, err = oauthClient.Do(req)

	case "POST":
		data, err := json.Marshal(cmd)
		if err != nil {
			return repMsg, err
		}
		reqData := bytes.NewBuffer(data)
		req, err := http.NewRequest("POST", url, reqData)
		if err != nil {
			return repMsg, err
		}
		req.Header.Set("Content-type", "application/json")
		response, err = oauthClient.Do(req)
	default:
		return repMsg, err
	}
	if response != nil {
		defer response.Body.Close()

		body, _ := ioutil.ReadAll(response.Body)
		if err := json.Unmarshal(body, &repMsg); err == nil {
			return repMsg, err
		}
	}
	if err != nil {
		return repMsg, err
	}

	return repMsg, err
}
func CheckAccessToken(auth_type, access_token string) (cakey string, errCode int64, result interface{}) {

	if len(Parameters.OauthServerAddr) == 0 {
		return "", Err.SUCCESS, ""
	}
	req := make(map[string]interface{})
	req["token"] = access_token
	req["auth_type"] = auth_type
	rep, err := OauthRequest("GET", req, Parameters.OauthServerAddr+"/"+access_token+"?auth_type="+auth_type)
	if err != nil {
		log.Error("Oauth timeout:", err)
		return "", Err.OAUTH_TIMEOUT, rep
	}
	if errcode, ok := rep["Error"].(float64); ok && errcode == 0 {
		result, ok := rep["Result"].(map[string]interface{})
		if !ok {
			return "", Err.INVALID_TOKEN, rep
		}
		if CAkey, ok := result["CaKey"].(string); ok {
			return CAkey, Err.SUCCESS, rep
		}
	}
	return "", Err.INVALID_TOKEN, rep
}
