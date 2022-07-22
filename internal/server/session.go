package server

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/myriadeinc/zircon/internal/stratum"
	"github.com/rs/zerolog/log"
)

const useDeadline = false

type StratumSession struct {
	sync.Mutex
	conn    *net.TCPConn
	ip      string
	minerId string
	service stratum.StratumService
}

func NewSession(ip string, connection *net.TCPConn, service stratum.StratumService) *StratumSession {
	return &StratumSession{
		ip:      ip,
		conn:    connection,
		service: service,
	}
}

func (s *StratumSession) closeSession() {
	log.Debug().Msg("Closing session")
	s.conn.Close()
	GetSessionHandler().removeSession(s)
	if r := recover(); r != nil {
		log.Error().Err(r.(error))
	}
}

func isEmptyRequest(rawRequest []byte) bool {
	s := strings.TrimSpace(string(rawRequest))
	return len(s) == 0
}

func (s *StratumSession) handleSession() {
	defer s.closeSession()

	timeoutDuration := 10 * time.Second
	clientReader := bufio.NewReader(s.conn)
	if useDeadline {
		s.conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	}
	for {
		// Could be long buffer, but at least we handle errors
		clientRequest, err := clientReader.ReadBytes('\n')
		if isEmptyRequest(clientRequest) {
			continue
		}
		log.Debug().Msgf("Raw read : %s", string(clientRequest))
		// Aggressively close session for any errors encountered
		switch err {
		case nil:
			requestErr := s.handleRequest(clientRequest)
			if requestErr != nil {
				log.Error().Msg("Disconnecting due to request process error")
				return
			}
			continue
		case io.EOF:
			log.Debug().Msgf("EOF on %s", s.ip)
			return
		default:
			log.Fatal().Msg("Error encountered while reading from connection")
			return
		}
	}
}

func (s *StratumSession) handleRequest(rawRequest []byte) error {
	var request JSONRpcReq
	jsonErr := json.Unmarshal(rawRequest, &request)
	if jsonErr != nil {
		log.Fatal().Err(jsonErr)
		return jsonErr
	}
	// if request.Method == "login" {
	// 	s.minerId = request.ParseMinerId()
	// }
	s.minerId = "00001111-1111-4222-8333-abc123456789"

	var handleErr error
	var response []byte
	var needNewJob bool

	log.Debug().Msg("Did I even reach here?")
	switch request.Method {
	case "login":
		job := s.service.HandleLogin(request.Id)
		response, handleErr = json.Marshal(job)
	case "submit":
		ok, newJob := s.service.HandleSubmit(request.Id)
		needNewJob = newJob
		response, handleErr = json.Marshal(ok)
	default:
		unknownMethod := map[string]string{
			"message": "unknownMethod",
		}
		response, handleErr = json.Marshal(unknownMethod)
	}
	if handleErr != nil {
		log.Fatal().Err(handleErr)
		log.Debug().Msg("oopsie!")
		return handleErr
	}
	log.Debug().Msgf("Request was %v", request)
	log.Debug().Msgf("Received request, response is %s", string(response))

	response = append(response, byte('\n'))
	_, writeErr := s.conn.Write(response)
	// Since we have two potential writes to socket, need to handle this concurrently in the future
	if needNewJob {
		go s.triggerNewJob()
	}
	return writeErr
}

func (s *StratumSession) triggerNewJob() {
	job := s.service.HandleNewJob()
	bytes, _ := json.Marshal(job)
	bytes = append(bytes, byte('\n'))
	s.conn.Write(bytes)
}
