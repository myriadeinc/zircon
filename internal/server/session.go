package server

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/myriadeinc/zircon/internal/cache"
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
	cache   cache.CacheService
}

func NewSession(ip string, connection *net.TCPConn, service stratum.StratumService, cache cache.CacheService) *StratumSession {
	return &StratumSession{
		ip:      ip,
		conn:    connection,
		service: service,
		cache:   cache,
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
		// Aggressively close session for any errors encountered
		switch err {
		case nil:
			requestErr := s.handleRequest(clientRequest)
			if requestErr != nil {
				log.Error().Err(requestErr).Msg("Disconnecting due to request process error")
				return
			}
			continue
		case io.EOF:
			log.Debug().Msgf("EOF on %s", s.ip)
			return
		default:
			log.Error().Msg("Error encountered while reading from connection")
			return
		}
	}
}

func (s *StratumSession) createLoginRequest(minerId string) (map[string]string, error) {
	template, err := s.cache.FetchTemplate()
	if err != nil {
		return nil, err
	}
	params := map[string]string{
		"miner":            minerId,
		"templateBlob":     template.BlockTemplateBlob,
		"templateHeight":   template.Height,
		"templateDiff":     template.Difficulty,
		"templateSeedhash": template.SeedHash,
	}

	return params, nil
}

func (s *StratumSession) handleRequest(rawRequest []byte) error {
	var request JSONRpcReq
	jsonErr := json.Unmarshal(rawRequest, &request)
	if jsonErr != nil {
		log.Error().Err(jsonErr).Msg("error unmarshaling request")
		return jsonErr
	}
	if request.Method == "login" {
		s.minerId = request.ParseMinerId()
	}

	var handleErr error
	var response []byte
	var pushNewJob bool

	log.Trace().Msgf("Raw request : %s", string(rawRequest))

	switch request.Method {
	case "login":
		params, err := s.createLoginRequest(s.minerId)
		if err != nil {
			response = genericErrorResponse(request.Id)
			break
		}
		job, err := s.service.HandleLoginWithTemplate(request.Id, params)
		if err != nil {
			response = genericErrorResponse(request.Id)
			break
		}
		response, handleErr = json.Marshal(job)
	case "submit":
		ok, err := s.service.HandleSubmit(request.Id, request.Params)
		if err != nil {
			response = genericErrorResponse(request.Id)
			break
		}
		response, handleErr = json.Marshal(ok)
		pushNewJob = true
	case "keepalived":
		keepalive := map[string]interface{}{
			"id":      request.Id,
			"jsonrpc": "2.0",
			"result": map[string]string{
				"status": "KEEPALIVED",
			},
		}
		response, _ = json.Marshal(keepalive)
	default:
		response = genericErrorResponse(request.Id)
	}
	if handleErr != nil {
		log.Error().Err(handleErr).Msg("Could not process incoming request")
		return handleErr
	}
	log.Trace().
		Str("request",
			string(rawRequest)).
		Str("response",
			string(response)).
		Msg("processed request")

	response = append(response, byte('\n'))
	_, writeErr := s.conn.Write(response)
	if writeErr != nil {
		return writeErr
	}
	if pushNewJob {
		return s.triggerNewJob()
	}
	return nil
}

func (s *StratumSession) triggerNewJob() error {
	params, err := s.createLoginRequest(s.minerId)
	if err != nil {
		return err
	}
	job, err := s.service.HandleNewJob(params)
	if err != nil {
		return err
	}
	bytes, _ := json.Marshal(job)

	bytes = append(bytes, byte('\n'))
	log.Trace().Msgf("Pushing new job %s", string(bytes))

	_, err = s.conn.Write(bytes)
	return err
}

func genericErrorResponse(id *json.RawMessage) []byte {
	response := map[string]interface{}{
		"id":      id,
		"jsonrpc": "2.0",
		"error": map[string]interface{}{
			"code":    -1,
			"message": "Internal server error",
		},
	}

	bytes, _ := json.Marshal(response)
	return bytes

}
