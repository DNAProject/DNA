package httpjsonrpc

import (
	. "github.com/DNAProject/DNA/common"
	. "github.com/DNAProject/DNA/config"
	"log"
	"net/http"
	"strconv"
)

func StartRPCServer() {
	Trace()
	http.HandleFunc("/", Handle)

	HandleFunc("getbestblockhash", getBestBlockHash)
	HandleFunc("getblock", getBlock)
	HandleFunc("getblockcount", getBlockCount)
	HandleFunc("getblockhash", getBlockHash)
	HandleFunc("getconnectioncount", getConnectionCount)
	HandleFunc("getrawmempool", getRawMemPool)
	HandleFunc("getrawtransaction", getRawTransaction)
	HandleFunc("submitblock", submitBlock)

	err := http.ListenAndServe(":"+strconv.Itoa(Parameters.HttpJsonPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
