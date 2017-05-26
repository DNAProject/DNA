package restful

import (
	. "DNA/common"
	. "DNA/common/config"
	"DNA/common/log"
	. "DNA/net/apiServer/common"
	"DNA/net/apiServer/errorcode"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

type handler func(map[string]interface{}) map[string]interface{}
type Handler struct {
	sync.RWMutex
	handler handler
	name    string
}
type restServer struct {
	postMap    map[string]Handler
	getMap     map[string]Handler
	HttpServer *http.Server
	Router     *Router
}

func InitRestServer() ApiServer {
	s := NewRestServer()
	return s
}
func NewRestServer() ApiServer {

	return &restServer{}
}

func (rt *restServer) Start() error {
	if Parameters.HttpRestPort == 0 {
		log.Fatal("Not configure HttpRestPort port ")
		return nil
	}
	rt.registryMethod()
	rt.Router = NewRouter()
	rt.initGetHandler()
	rt.initPostHandler()

	rest := false
	if !rest {
		err := http.ListenAndServe(":"+strconv.Itoa(Parameters.HttpRestPort), rt.Router)
		if err != nil {
			log.Fatal("ListenAndServe: ", err.Error())
			return err
		}
		return nil
	}

	CAPath := ""
	ctPool := x509.NewCertPool()
	byrCrtData, err := ioutil.ReadFile(CAPath)
	if err != nil {
		log.Error("ReadFile rootCA.crt Error:", err)
		return err
	}
	ctPool.AppendCertsFromPEM(byrCrtData)
	rt.HttpServer = &http.Server{
		Addr: ":" + strconv.Itoa(Parameters.HttpRestPort),
		TLSConfig: &tls.Config{
			ClientCAs:  ctPool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}
	err = rt.HttpServer.ListenAndServeTLS("cert.pem", "key.pem")
	if err != nil {
		return err
	}

	return nil
}

func (rt *restServer) registryMethod() {

	getMethodMap := map[string]Handler{
		"/api/v1/node/info":             {handler: GetNodeInfo, name: "nodeinfo"},
		"/api/v1/node/count":            {handler: GetNodeCount, name: "nodecount"},
		"/api/v1/block/height":          {handler: GetBlockHeight, name: "getblockheight"},
		"/api/v1/block/info/height":     {handler: GetBlockByHeight, name: "getblockbyheight"},
		"/api/v1/block/info/hash":       {handler: GetBlockByHash, name: "getblockbyhash"},
		"/api/v1/asset":                 {handler: GetAssetByHash, name: "getasset"},
		"/api/v1/record":                {handler: GetRecordByHash, name: "getrecord"},
		"/api/v1/transaction":           {handler: GetTransactionByHash, name: "gettransaction"},
		"/api/v1/transaction/inmempool": {handler: GetTransactionsInMempool, name: "transactionsinmempool"},
		"/api/v1/smartcode/info":        {handler: GetSmartcodeInfo, name: "smartcodeinfo"},
		"/api/v1/test":                  {handler: Test},
	}

	postMethodMap := map[string]Handler{
		"/api/v1/block":             {handler: PostBlock, name: "postblock"},
		"/api/v1/asset/transfer":    {handler: PostAssetTransfer, name: "assettransfer"},
		"/api/v1/asset/issue":       {handler: PostAssetIssue, name: "assetissue"},
		"/api/v1/asset/registry":    {handler: PostAssetRegistry, name: "assetregistry"},
		"/api/v1/record":            {handler: PostRecord, name: "postrecord"},
		"/api/v1/smartcode/publish": {handler: PostSmartcodePublish, name: "smartcodepublish"},
		"/api/v1/smartcode/invoke":  {handler: PostSmartcodeInvoke, name: "smartcodeinvoke"},
	}
	rt.postMap = postMethodMap
	rt.getMap = getMethodMap
}
func (rt *restServer) initGetHandler() {

	for k, _ := range rt.getMap {
		rt.Router.Get(k, func(w http.ResponseWriter, r *http.Request) {

			var reqMsg = make(map[string]interface{})

			if h, ok := rt.getMap[r.URL.Path]; ok {
				h.Lock()
				defer h.Unlock()

				reqMsg["id"] = r.FormValue("id")
				reqMsg["height"] = r.FormValue("height")
				reqMsg["hash"] = r.FormValue("hash")
				resp := h.handler(reqMsg)

				resp["action"] = h.name + "-ack"
				data, err := json.Marshal(resp)
				if err != nil {
					log.Fatal("HTTP Handle - json.Marshal: %v", err)
					return
				}
				w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
				w.Header().Set("content-type", "application/json")
				w.Write([]byte(data))

			}
		})
	}
}
func (rt *restServer) initPostHandler() {
	for k, _ := range rt.postMap {
		rt.Router.Post(k, func(w http.ResponseWriter, r *http.Request) {

			body, _ := ioutil.ReadAll(r.Body)
			var reqMsg = make(map[string]interface{})
			var data []byte
			if h, ok := rt.postMap[r.URL.Path]; ok {

				if err := json.Unmarshal(body, &reqMsg); err == nil {
					h.Lock()
					defer h.Unlock()

					resp := h.handler(reqMsg)
					resp["id"] = r.FormValue("id")
					resp["action"] = h.name + "-ack"
					data, err = json.Marshal(resp)
					if err != nil {
						log.Fatal("HTTP Handle - json.Marshal: %v", err)
						return
					}
					//w.Header().Set("Access-Control-Allow-Origin", "*")//

				} else {
					var resp = make(map[string]interface{})
					resp["id"] = r.FormValue("id")
					resp["action"] = h.name + "-ack"
					resp["error"] = errorcode.ILLEGAL_DATAFORMAT
					data, err = json.Marshal(resp)
					if err != nil {
						log.Fatal("HTTP Handle - json.Marshal: %v", err)
						return
					}
				}
				w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
				w.Header().Set("content-type", "application/json")
				w.Write([]byte(data))
			}
		})
	}
	//
	for k, _ := range rt.postMap {
		rt.Router.Options(k, func(w http.ResponseWriter, r *http.Request) {
		})
	}

}
func (rt *restServer) Stop() {
	//TODO
}

func (rt *restServer) Push(txHash Uint256, errcode interface{}, result interface{}) {

}
