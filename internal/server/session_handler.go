package server

import (
	"sync"

	"github.com/rs/zerolog/log"
)

type SessionHandler struct {
	sync.RWMutex
	Sessions map[*StratumSession]struct{}
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

func (s *SessionHandler) BroadcastNewJobs() {
	log.Info().Msg("broadcasting new jobs")
	badSessions := []*StratumSession{}
	s.Lock()
	defer s.Unlock()
	var wg sync.WaitGroup
	wg.Add(len(s.Sessions))
	for session := range s.Sessions {
		go func(sess *StratumSession) {
			log.Trace().Msgf("trigger new job for %s", sess.ip)
			err := sess.triggerNewJob()
			if err != nil {
				log.Error().Err(err).Msgf("Could not push msg to session %s", sess.ip)
				badSessions = append(badSessions, sess)
			}
			wg.Done()
		}(session)
	}
	wg.Wait()
	for _, session := range badSessions {
		s.removeSession(session)
	}
}

func (s *SessionHandler) removeSession(session *StratumSession) {
	s.Lock()
	defer s.Unlock()
	delete(s.Sessions, session)
}
func (s *SessionHandler) addSession(session *StratumSession) {
	s.Lock()
	defer s.Unlock()
	s.Sessions[session] = struct{}{}
}
