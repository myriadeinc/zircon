package server

import (
	"strconv"
	"time"

	"github.com/myriadeinc/zircon/internal/models"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var samplog = log.Sample(&zerolog.BasicSampler{N: 10})

func (s *PoolServer) pollForever() {
	log.Info().Uint64("blockheight", s.blockHeight).Msg("Starting poller")
	// samplog := log.Sample(&zerolog.BasicSampler{N: 10})

	for {
		template, err := s.NodeClient.GetValidBlockTemplate()
		// Case of bad startup
		if err != nil {
			log.Error().Err(err).Msg("retry in 15 seconds")
			time.Sleep(15 * time.Second)
			continue
		}
		height, err := strconv.ParseUint(template.Height, 10, 64)
		if err != nil {
			log.Error().Err(err).Msg("retry in 15 seconds")
			time.Sleep(15 * time.Second)
			continue
		}

		err = s.Cache.SaveNewTemplate(*template)
		if err != nil {
			log.Error().Err(err).Msg("could not save newtemplate retry in 15 seconds")
			time.Sleep(15 * time.Second)
			continue
		}
		log.Trace().Msg("Saved new blocktemplate")
		if height > s.blockHeight {
			samplog.Info().Uint64("height", height).Uint64("prevHeight", s.blockHeight).Msg("new blockheight broadcast jobs")
			s.blockHeight = height
			// Fire and forget
			go s.BroadcastNewJobs(*template)
		}

		time.Sleep(15 * time.Second)

	}

}

func (s *PoolServer) BroadcastNewJobs(template models.StrictTemplate) {
	s.SessionLock.RLock()
	numSessions := len(s.Sessions)
	samplog.Info().Int("nsessions", numSessions).Msg("broadcasting new jobs")
	activeSessions := 0
	badSessions := []*StratumSession{}
	for session := range s.Sessions {

		log.Trace().Str("ip", session.ip).Msg("trigger new job")
		err := s.triggerNewJob(session, template)
		if err != nil {
			log.Error().Err(err).Str("ip", session.ip).Msg("Could not push msg to session")
			badSessions = append(badSessions, session)
		} else {
			activeSessions++
		}

	}
	s.SessionLock.RUnlock()

	for _, session := range badSessions {
		s.removeSession(session)
	}

	samplog.Info().Int("nsessions", activeSessions).Msg("Sucessfully broadcasted jobs")

}
