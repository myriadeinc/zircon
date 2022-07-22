package server

import (
	"sync"
)

type SessionHandler struct {
	SessionsMu sync.RWMutex
	Sessions   map[*StratumSession]struct{}
}

var once sync.Once
var sessionHandler *SessionHandler

func GetSessionHandler() *SessionHandler {
	if sessionHandler == nil {
		once.Do(
			func() {
				sessionMap := make(map[*StratumSession]struct{})
				sessionHandler = &SessionHandler{
					Sessions: sessionMap,
				}
			})
	}
	return sessionHandler
}

func (s *SessionHandler) removeSession(session *StratumSession) {
	s.SessionsMu.Lock()
	defer s.SessionsMu.Unlock()
	delete(s.Sessions, session)
}
func (s *SessionHandler) addSession(session *StratumSession) {
	s.SessionsMu.Lock()
	defer s.SessionsMu.Unlock()
	s.Sessions[session] = struct{}{}
}
