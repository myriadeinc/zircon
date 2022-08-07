package server

import (
	"bufio"
	"encoding/json"
	"io"

	// "github.com/myriadeinc/zircon/internal/models"
	"github.com/myriadeinc/zircon/internal/models"
	"github.com/rs/zerolog/log"
)

func (s *PoolServer) handlePoolSession(session *StratumSession) {
	clientReader := bufio.NewReader(session.conn)
	defer s.removeSession(session)
	for {
		// Could be long buffer, but at least we handle errors
		clientRequest, err := clientReader.ReadBytes('\n')
		if isEmptyRequest(clientRequest) {
			continue
		}
		// Aggressively close session for any errors encountered
		switch err {
		case nil:
			requestErr := s.handleRequest(session, clientRequest)
			if requestErr != nil {
				log.Error().Err(requestErr).Msg("request process error")
				return
			}
			continue
		case io.EOF:
			log.Error().Err(err).Msgf("EOF on %s", session.ip)
			return
		default:
			log.Error().Err(err).Msg("Error encountered while reading from connection")
			return
		}
	}

}

func (s *PoolServer) handleRequest(session *StratumSession, rawRequest []byte) error {
	var request JSONRpcReq
	jsonErr := json.Unmarshal(rawRequest, &request)
	if jsonErr != nil {
		log.Error().Err(jsonErr).Msg("error unmarshaling request")
		return jsonErr
	}
	if request.Method == "login" {
		session.minerId = request.ParseMinerId()
	}

	var handleErr error
	var response []byte
	var pushNewJob bool

	log.Trace().Msgf("Raw request : %s", string(rawRequest))

	switch request.Method {
	case "login":
		template, err := s.Cache.FetchTemplate()
		if err != nil {
			response = genericErrorResponse(request.Id)
			break
		}
		params := template.ToLoginRequest(session.minerId)
		job, err := s.Stratum.HandleLoginWithTemplate(request.Id, params)
		if err != nil {
			response = genericErrorResponse(request.Id)
			break
		}
		response, handleErr = json.Marshal(job)
	case "submit":
		ok, err := s.Stratum.HandleSubmit(request.Id, request.Params)
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
	_, writeErr := session.conn.Write(response)
	if writeErr != nil {
		return writeErr
	}
	if pushNewJob {
		template, err := s.Cache.FetchTemplate()
		if err != nil {
			return err
		}
		return s.triggerNewJob(session, *template)
	}
	return nil
}

func (s *PoolServer) triggerNewJob(session *StratumSession, template models.StrictTemplate) error {
	params := template.ToLoginRequest(session.minerId)
	job, err := s.Stratum.HandleNewJob(params)
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(job)
	if err != nil {
		return err
	}

	bytes = append(bytes, byte('\n'))
	log.Trace().Msgf("Pushing new job %s", string(bytes))

	_, err = session.conn.Write(bytes)
	return err
}

func (s *PoolServer) addSession(session *StratumSession) {
	s.SessionLock.Lock()
	defer s.SessionLock.Unlock()
	s.Sessions[session] = struct{}{}
}

func (s *PoolServer) removeSession(session *StratumSession) {
	s.SessionLock.Lock()
	defer s.SessionLock.Unlock()
	delete(s.Sessions, session)
	session.conn.Close()
}
