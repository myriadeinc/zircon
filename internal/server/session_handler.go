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
	badSessions := []*StratumSession{}
	s.Lock()
	defer s.Unlock()
	var wg sync.WaitGroup
	numSessions := len(s.Sessions)
	wg.Add(numSessions)
	log.Info().Int("nsessions", numSessions).Msg("broadcasting new jobs")
	activeSessions := 0
	for session := range s.Sessions {
		go func(sess *StratumSession) {
			log.Trace().Msgf("trigger new job for %s", sess.ip)
			err := sess.triggerNewJob()
			if err != nil {
				log.Error().Err(err).Str("ip", sess.ip).Msg("Could not push msg to session")
				badSessions = append(badSessions, sess)
			} else {
				activeSessions++
			}
			wg.Done()
		}(session)
	}
	wg.Wait()
	for _, session := range badSessions {
		s.removeSession(session)
	}
	log.Info().Int("nsessions", activeSessions).Msg("Sucessfully broadcasted jobs")
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
