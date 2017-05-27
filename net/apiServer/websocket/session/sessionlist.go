package session

import (
	"sync"
)

type SessionList struct {
	sync.RWMutex
	mapOnlineList map[string]*Session //key is SessionId
}

func NewSessionList() *SessionList {
	return &SessionList{
		mapOnlineList: make(map[string]*Session),
	}
}

func (sl *SessionList) addOnlineSession(session *Session) {
	if session.GetSessionId() != "" {
		sl.Lock()
		defer sl.Unlock()
		sl.mapOnlineList[session.GetSessionId()] = session
	}
}

func (sl *SessionList) removeSession(iSession *Session) (err error) {

	if iSession.GetSessionId() != "" {
		sl.Lock()
		defer sl.Unlock()
		delete(sl.mapOnlineList, iSession.GetSessionId())
	}
	return err
}

func (sl *SessionList) removeSessionById(sSessionId string) (err error) {

	if sSessionId != "" {
		sl.Lock()
		defer sl.Unlock()
		delete(sl.mapOnlineList, sSessionId)
	}
	return err
}

func (sl *SessionList) GetSessionById(sSessionId string) *Session {
	sl.RLock()
	defer sl.RUnlock()
	session, bOk := sl.mapOnlineList[sSessionId]
	if bOk {
		return session
	}
	return nil

}
func (sl *SessionList) GetSessionList() map[string]*Session {
	return sl.mapOnlineList
}
