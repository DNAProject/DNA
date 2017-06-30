package websocket

import (
	. "DNA/common/config"
	"DNA/common/log"
	. "DNA/net/httprestful/common"
	Err "DNA/net/httprestful/error"
	. "DNA/net/httpwebsocket/session"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type handler func(map[string]interface{}) map[string]interface{}
type Handler struct {
	handler  handler
	pushFlag bool
}

type WsServer struct {
	sync.RWMutex
	Upgrader         websocket.Upgrader
	listener         net.Listener
	server           *http.Server
	SessionList      *SessionList
	ActionMap        map[string]Handler
	TxHashMap        map[string]string //key: txHash   value:sessionid
	checkAccessToken func(auth_type, access_token string) (string, int64, interface{})
}

func InitWsServer(checkAccessToken func(string, string) (string, int64, interface{})) *WsServer {
	ws := &WsServer{
		Upgrader:    websocket.Upgrader{},
		SessionList: NewSessionList(),
		TxHashMap:   make(map[string]string),
	}
	ws.checkAccessToken = checkAccessToken
	return ws
}

func (ws *WsServer) Start() error {
	if Parameters.HttpWsPort == 0 {
		log.Error("Not configure HttpWsPort port ")
		return nil
	}
	ws.registryMethod()
	ws.Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	tlsFlag := false
	if tlsFlag || Parameters.HttpWsPort%1000 == 443 {
		var err error
		ws.listener, err = ws.initTlsListen()
		if err != nil {
			log.Error("Https Cert: ", err.Error())
			return err
		}
	} else {
		var err error
		ws.listener, err = net.Listen("tcp", ":"+strconv.Itoa(Parameters.HttpWsPort))
		if err != nil {
			log.Fatal("net.Listen: ", err.Error())
			return err
		}
	}
	ws.server = &http.Server{Handler: http.HandlerFunc(ws.webSocketHandler)}
	err := ws.server.Serve(ws.listener)

	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
		return err
	}
	return nil

}

func (ws *WsServer) registryMethod() {
	gettxhashmap := func(cmd map[string]interface{}) map[string]interface{} {
		resp := ResponsePack(Err.SUCCESS)
		ws.Lock()
		defer ws.Unlock()
		resp["Result"] = ws.TxHashMap
		return resp
	}
	sendRawTransaction := func(cmd map[string]interface{}) map[string]interface{} {
		resp := SendRawTransaction(cmd)
		if userid, ok := resp["Userid"].(string); ok && len(userid) > 0 {
			if result, ok := resp["Result"].(string); ok {
				ws.SetTxHashMap(result, userid)
			}
			delete(resp, "Userid")
		}
		return resp
	}
	heartbeat := func(cmd map[string]interface{}) map[string]interface{} {
		resp := ResponsePack(Err.SUCCESS)
		resp["Action"] = "heartbeat"
		resp["Result"] = cmd["Userid"]
		return resp
	}
	sendtest := func(cmd map[string]interface{}) map[string]interface{} {
		go func() {
			time.Sleep(time.Second * 5)
			resp := ResponsePack(Err.SUCCESS)
			resp["Action"] = "pushresult"
			ws.PushTxResult(cmd["Userid"].(string), resp)
		}()
		return heartbeat(cmd)
	}
	getsessionlist := func(cmd map[string]interface{}) map[string]interface{} {
		resp := ResponsePack(Err.SUCCESS)
		resp["Action"] = "getsessionlist"
		resp["Result"] = ws.SessionList.GetSessionList()
		return resp
	}
	actionMap := map[string]Handler{
		"getconnectioncount": {handler: GetConnectionCount},
		"getblockbyheight":   {handler: GetBlockByHeight},
		"getblockbyhash":     {handler: GetBlockByHash},
		"getblockheight":     {handler: GetBlockHeight},
		"gettransaction":     {handler: GetTransactionByHash},
		"getasset":           {handler: GetAssetByHash},
		"getunspendoutput":   {handler: GetUnspendOutput},

		"sendrawtransaction": {handler: sendRawTransaction},
		"sendrecord":         {handler: SendRecorByTransferTransaction},
		"heartbeat":          {handler: heartbeat},

		"sendtest": {handler: sendtest, pushFlag: true},

		"gettxhashmap":   {handler: gettxhashmap},
		"getsessionlist": {handler: getsessionlist},
	}
	ws.ActionMap = actionMap
}

func (ws *WsServer) Stop() {
	if ws.server != nil {
		ws.server.Shutdown(context.Background())
		log.Error("Close websocket ")
	}
}
func (ws *WsServer) Restart() {
	go func() {
		time.Sleep(time.Second)
		ws.Stop()
		time.Sleep(time.Second)
		go ws.Start()
	}()
}

//webSocketHandler
func (ws *WsServer) webSocketHandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := ws.Upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Error("ListenAndServe: ", err.Error())
		return
	}
	nsSession, err := ws.SessionList.NewSession(wsConn)
	if err != nil {
		log.Error("websocket NewSession fail.")
		return
	}

	defer func() {
		ws.deleteTxHashs(nsSession.GetSessionId())
		ws.SessionList.CloseSession(nsSession)
		if err := recover(); err != nil {
			log.Fatal("websocket recover.")
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
			}
			continue
		} else {
			ws.deleteTxHashs(nsSession.GetSessionId())
			ws.SessionList.CloseSession(nsSession)
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
func (ws *WsServer) IsValidMsg(reqMsg map[string]interface{}) bool {
	if _, ok := reqMsg["Hash"].(string); !ok && reqMsg["Hash"] != nil {
		return false
	}
	if _, ok := reqMsg["Addr"].(string); !ok && reqMsg["Addr"] != nil {
		return false
	}
	if _, ok := reqMsg["Assetid"].(string); !ok && reqMsg["Assetid"] != nil {
		return false
	}
	return true
}
func (ws *WsServer) OnDataHandle(curSession *Session, bysMsg []byte, r *http.Request) bool {

	var req = make(map[string]interface{})

	if err := json.Unmarshal(bysMsg, &req); err != nil {
		resp := ResponsePack(Err.ILLEGAL_DATAFORMAT)
		ws.response(curSession.GetSessionId(), resp)
		log.Error("websocket OnDataHandle ILLEGAL_DATAFORMAT")
		return false
	}
	actionName, ok := req["Action"].(string)
	if !ok {
		resp := ResponsePack(Err.INVALID_METHOD)
		ws.response(curSession.GetSessionId(), resp)
		return true
	}
	action, ok := ws.ActionMap[actionName]
	if !ok {
		resp := ResponsePack(Err.INVALID_METHOD)
		ws.response(curSession.GetSessionId(), resp)
		return true
	}
	if !ws.IsValidMsg(req) {
		resp := ResponsePack(Err.INVALID_PARAMS)
		ws.response(curSession.GetSessionId(), resp)
		return true
	}
	if height, ok := req["Height"].(float64); ok {
		req["Height"] = strconv.FormatInt(int64(height), 10)
	}
	if raw, ok := req["Raw"].(float64); ok {
		req["Raw"] = strconv.FormatInt(int64(raw), 10)
	}
	auth_type, ok := req["auth_type"].(string)
	if !ok {
		auth_type = ""
	}
	access_token, ok := req["access_token"].(string)
	if !ok {
		access_token = ""
	}
	if actionName != "heartbeat" {
		CAkey, errCode, result := ws.checkAccessToken(auth_type, access_token)
		if errCode > 0 {
			resp := ResponsePack(errCode)
			resp["Result"] = result
			ws.response(curSession.GetSessionId(), resp)
			return true
		}
		req["CAkey"] = CAkey
	}
	req["Userid"] = curSession.GetSessionId()
	resp := action.handler(req)
	resp["Action"] = actionName
	if txHash, ok := resp["Result"].(string); ok && action.pushFlag {
		ws.Lock()
		defer ws.Unlock()
		ws.TxHashMap[txHash] = curSession.GetSessionId()
	}
	ws.response(curSession.GetSessionId(), resp)

	return true
}
func (ws *WsServer) SetTxHashMap(txhash string, sessionid string) {
	ws.Lock()
	defer ws.Unlock()
	ws.TxHashMap[txhash] = sessionid
}
func (ws *WsServer) deleteTxHashs(sSessionId string) {
	ws.Lock()
	defer ws.Unlock()
	for k, v := range ws.TxHashMap {
		if v == sSessionId {
			delete(ws.TxHashMap, k)
		}
	}
}
func (ws *WsServer) response(sSessionId string, resp map[string]interface{}) {
	resp["Desc"] = Err.ErrMap[resp["Error"].(int64)]
	data, err := json.Marshal(resp)
	if err != nil {
		log.Error("Websocket Handle - json.Marshal: %v", err)
		return
	}
	ws.send(sSessionId, data)
}
func (ws *WsServer) PushTxResult(txHashStr string, resp map[string]interface{}) {
	ws.Lock()
	defer ws.Unlock()
	sSessionId := ws.TxHashMap[txHashStr]
	delete(ws.TxHashMap, txHashStr)
	if len(sSessionId) > 0 {
		ws.response(sSessionId, resp)
	}
}
func (ws *WsServer) PushResult(resp map[string]interface{}) {
	resp["Desc"] = Err.ErrMap[resp["Error"].(int64)]
	data, err := json.Marshal(resp)
	if err != nil {
		log.Error("Websocket Handle - json.Marshal: %v", err)
		return
	}
	ws.broadcast(data)
}
func (ws *WsServer) send(sSessionId string, data []byte) error {
	session := ws.SessionList.GetSessionById(sSessionId)
	if session == nil {
		return errors.New("SessionId Not Exist:" + sSessionId)
	}
	return session.Send(data)
}
func (ws *WsServer) broadcast(data []byte) error {
	for _, v := range ws.SessionList.GetSessionList() {
		if session := ws.SessionList.GetSessionById(v); session != nil {
			session.Send(data)
		}
	}
	return nil
}

func (ws *WsServer) initTlsListen() (net.Listener, error) {

	CertPath := Parameters.RestCertPath
	KeyPath := Parameters.RestKeyPath

	// load cert
	cert, err := tls.LoadX509KeyPair(CertPath, KeyPath)
	if err != nil {
		log.Error("load keys fail", err)
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	log.Info("TLS listen port is ", strconv.Itoa(Parameters.HttpWsPort))
	listener, err := tls.Listen("tcp", ":"+strconv.Itoa(Parameters.HttpWsPort), tlsConfig)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return listener, nil
}
