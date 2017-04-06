package httpjsonrpc

import (
	"DNA/common/log"
	. "DNA/config"
	"net/http"
	"strconv"
)

func StartRPCServer() {
	log.Trace()
	http.HandleFunc("/", Handle)

	HandleFunc("getbestblockhash", getBestBlockHash)
	HandleFunc("getblock", getBlock)
	HandleFunc("getTxn", getTxn)
	HandleFunc("getAddrTxn", getAddrTxn)
	HandleFunc("getblockcount", getBlockCount)
	HandleFunc("getblockhash", getBlockHash)
	HandleFunc("getconnectioncount", getConnectionCount)
	HandleFunc("getrawmempool", getRawMemPool)
	HandleFunc("getrawtransaction", getRawTransaction)
	HandleFunc("sendRawTransaction", sendRawTransaction)
	HandleFunc("submitblock", submitBlock)

	err := http.ListenAndServe(":"+strconv.Itoa(Parameters.HttpJsonPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
