package session

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/pborman/uuid"
	"sync"
	"time"
)

type Session struct {
	sync.Mutex
	mSessionList *SessionList
	mConnection  *websocket.Conn
	nLastActive  int64
	sSessionId   string
}

func (s *Session) GetSessionId() string {
	return s.sSessionId
}

func NewSession(sessionList *SessionList, wsConn *websocket.Conn) (hispSession *Session, err error) {
	sSessionId := uuid.NewUUID().String()
	hispSession = &Session{
		mConnection:  wsConn,
		nLastActive:  time.Now().Unix(),
		sSessionId:   sSessionId,
		mSessionList: sessionList,
	}
	sessionList.addOnlineSession(hispSession)
	return hispSession, err
}

func (s *Session) Close() {
	if s.mConnection != nil {
		s.mSessionList.removeSession(s)
		s.mConnection.Close()
		s.mConnection = nil
	}
	s.sSessionId = ""
}

func (s *Session) UpdateActiveTime() {
	s.nLastActive = time.Now().Unix()
}

func (s *Session) Send(data []byte) error {
	if s.mConnection == nil {
		return errors.New("WebSocket is null")
	}
	//https://godoc.org/github.com/gorilla/websocket
	s.Lock()
	defer s.Unlock()
	return s.mConnection.WriteMessage(websocket.TextMessage, data)
}

func (s *Session) SessionTimeoverCheck() bool {

	nCurTime := time.Now().Unix()
	if nCurTime-s.nLastActive > 300 { //5 mins
		return true
	}
	return false
}
