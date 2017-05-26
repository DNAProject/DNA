package websocket

import (
	. "DNA/common"
	. "DNA/common/config"
	"DNA/common/log"
	. "DNA/net/apiServer/common"
	"DNA/net/apiServer/errorcode"
	. "DNA/net/apiServer/websocket/session"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type handler func(map[string]interface{}) map[string]interface{}
type wsHandler struct {
	handler handler
	push    bool
}

type wsServer struct {
	sync.RWMutex
	Upgrader    websocket.Upgrader
	HttpServer  *http.Server
	SessionList *SessionList
	MethodMap   map[string]wsHandler
	TxMap       map[Uint256]string //key: txHash   value:sessionid
}

func InitWsServer() ApiServer {
	s := NewWsServer()
	return s
}
func NewWsServer() *wsServer {

	return &wsServer{
		Upgrader:    websocket.Upgrader{},
		SessionList: NewSessionList(),
		TxMap:       make(map[Uint256]string),
	}
}

func (ws *wsServer) registryMethod() {
	methodMap := map[string]wsHandler{
		"nodeinfo":              {GetNodeInfo, false},
		"nodecount":             {GetNodeCount, false},
		"getblockheight":        {GetBlockHeight, false},
		"getblockbyhash":        {GetBlockByHash, false},
		"getblockbyheight":      {GetBlockByHeight, false},
		"postblock":             {PostBlock, false},
		"assetinfo":             {GetAssetByHash, false},
		"assettransfer":         {PostAssetTransfer, false},
		"assetissue":            {PostAssetIssue, false},
		"assetregistry":         {PostAssetRegistry, false},
		"postrecord":            {PostRecord, false},
		"getrecord":             {GetRecordByHash, false},
		"gettransaction":        {GetTransactionByHash, false},
		"transactionsinmempool": {GetTransactionsInMempool, false},
		"smartcodepublish":      {PostSmartcodePublish, true},
		"smartcodeinvoke":       {PostSmartcodeInvoke, true},
		"smartcodeinfo":         {GetSmartcodeInfo, false},
		"heartbeat":             {WsHeartbeat, false},
		"test":                  {Test, false},
	}
	ws.MethodMap = methodMap
}

func (ws *wsServer) Start() error {
	if Parameters.HttpWsPort == 0 {
		log.Fatal("Not configure HttpWsPort port ")
		return nil
	}

	ws.registryMethod()
	http.HandleFunc("/ws", ws.webSocketHandler)

	ws.Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	wss := false
	if !wss {
		err := http.ListenAndServe(":"+strconv.Itoa(Parameters.HttpWsPort), nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err.Error())
			return err
		}
		return nil
	}
	//TODO test tls
	CAPath := ""
	ctPool := x509.NewCertPool()
	byrCrtData, err := ioutil.ReadFile(CAPath)
	if err != nil {
		log.Error("ReadFile rootCA.crt Error:", err)
		return err
	}
	ctPool.AppendCertsFromPEM(byrCrtData)
	ws.HttpServer = &http.Server{
		Addr:    ":" + strconv.Itoa(Parameters.HttpWsPort),
		Handler: http.HandlerFunc(ws.webSocketHandler),
		TLSConfig: &tls.Config{
			ClientCAs:  ctPool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}
	err = ws.HttpServer.ListenAndServeTLS("cert.pem", "key.pem")
	if err != nil {

	}
	return err

}
func (ws *wsServer) Stop() {
	//TODO
}

//webSocketHandler
func (ws *wsServer) webSocketHandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := ws.Upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Error("ListenAndServe: ", err.Error())
		return
	}
	nsSession, err := NewSession(ws.SessionList, wsConn)
	if err != nil {
		nsSession.Close()
		return
	}

	defer func() {
		nsSession.Close()
		if err := recover(); err != nil {
		}
	}()

	for {
		//Set Read Deadline
		err = wsConn.SetReadDeadline(time.Now().Add(time.Second * 30))
		if err != nil {
			return
		}
		msgType, bysMsg, err := wsConn.ReadMessage()
		if err == nil && msgType == websocket.TextMessage {
			if ws.OnDataHandle(nsSession, bysMsg, r) {
				nsSession.UpdateActiveTime() //Update Active Time
			} else {

			}
			continue
		} else {
			nsSession.Close()
			return
		}

		//error and timeoutcheck
		e, ok := err.(net.Error)
		if !ok || !e.Timeout() {
			return
		} else if nsSession.SessionTimeoverCheck() {
			return
		}
	}
}

func (ws *wsServer) OnDataHandle(curSession *Session, bysMsg []byte, r *http.Request) bool {

	var reqMsg = make(map[string]interface{})

	if err := json.Unmarshal(bysMsg, &reqMsg); err != nil {
		log.Error("OnDataHandle:")
		ws.response(curSession.GetSessionId(), errorcode.ILLEGAL_DATAFORMAT, 0, 0, "ack")
		return true
	}

	if reqMsg["action"] == nil {
		ws.response(curSession.GetSessionId(), errorcode.INVALID_METHOD, 0, 0, "ack")
		return true
	}
	actionName := reqMsg["action"].(string)
	h, ok := ws.MethodMap[actionName]
	if !ok {
		ws.response(curSession.GetSessionId(), errorcode.INVALID_METHOD, 0, 0, "ack")
		return true
	}
	if !h.push {
		repMsg := h.handler(reqMsg)
		ws.response(curSession.GetSessionId(), errorcode.SUCCESS, repMsg["result"], repMsg["id"], actionName+"-ack")
		return true
	}
	msgTxHash := Uint256{} //TODO getTxHash(reqMsg)

	ws.Lock()
	defer ws.Unlock()

	if ws.TxMap[msgTxHash] == "" {
		repMsg := h.handler(reqMsg)
		ws.TxMap[msgTxHash] = curSession.GetSessionId()
		ws.response(curSession.GetSessionId(), errorcode.SUCCESS, repMsg["result"], repMsg["id"], actionName+"-ack")
	} else {
		ws.response(curSession.GetSessionId(), errorcode.PROCESSING_TX, "Processing transaction", reqMsg["id"], actionName+"-ack")
	}

	return true
}
func (ws *wsServer) responsePack(errcode interface{}, result interface{}, id interface{}, action interface{}) []byte {
	resp := map[string]interface{}{
		"action": action,
		"error":  errcode,
		"result": result,
		"id":     id,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		log.Fatal("Websocket Handle - json.Marshal: %v", err)
	}
	return data
}
func (ws *wsServer) response(sSessionId string, errcode interface{}, result interface{}, id interface{}, action interface{}) {

	data := ws.responsePack(errcode, result, id, action)
	err := ws.send(sSessionId, data)
	//err := ws.Broadcast(data)
	if err != nil {
	}
}
func (ws *wsServer) Push(txHash Uint256, errcode interface{}, result interface{}) {

	if ws.TxMap[txHash] == "" {
		return
	}
	data := ws.responsePack(errcode, result, "", "push-ack")
	//err = this.Send(ws.TxMap[txHash], data)
	err := ws.broadcast(data)
	if err != nil {
	}
	delete(ws.TxMap, txHash)
}
func (ws *wsServer) send(sSessionId string, data []byte) error {

	session := ws.SessionList.GetSessionById(sSessionId)
	if session == nil {
		return errors.New("SessionId Not Exist:" + sSessionId)
	}
	return session.Send(data)
}
func (ws *wsServer) broadcast(data []byte) error {
	for _, session := range ws.SessionList.GetSessionList() {
		session.Send(data)
	}
	return nil
}
func (ws *wsServer) closeSession(sSessionId string) error {
	session := ws.SessionList.GetSessionById(sSessionId)
	if session == nil {
		return errors.New("SessionId Not Exist:" + sSessionId)
	}
	session.Close()
	return nil
}
