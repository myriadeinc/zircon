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
}

func NewSession(ip string, connection *net.TCPConn) *StratumSession {
	return &StratumSession{
		ip:   ip,
		conn: connection,
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
		// log.Debug().Msgf("Raw read : %s", string(clientRequest))
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
		return jsonErr
	}
	if request.Method == "login" {
		s.minerId = request.ParseMinerId()
	}
	needNewJob, response, stratumErr := request.GetStratumResponse(s.minerId)
	if stratumErr != nil {
		return stratumErr
	}

	response = append(response, byte('\n'))
	_, writeErr := s.conn.Write(response)
	// Since we have two potential writes to socket, need to handle this concurrently in the future
	if needNewJob {
		go s.triggerNewJob()
	}
	return writeErr
}

func (s *StratumSession) triggerNewJob() {
	job := stratum.GetDummyJobWrapper()
	bytes, _ := json.Marshal(job)
	bytes = append(bytes, byte('\n'))
	s.conn.Write(bytes)
}

// func (st *StratumSession) sendResult(id *json.RawMessage, result interface{}, err error) error {
// 	if err != nil {
// 		message := JSONRpcResp{Id: id, Version: "2.0", Error: err, Result: nil}
// 		return st.enc.Encode(&message)
// 	}
// 	st.Lock()
// 	defer st.Unlock()
// 	message := JSONRpcResp{Id: id, Version: "2.0", Error: nil, Result: result}
// 	return st.enc.Encode(&message)
// }

// // Push notification
// func (st *StratumSession) sendJob(result interface{}) error {
// 	st.Lock()
// 	defer st.Unlock()
// 	message := JSONRpcPushMessage{Version: "2.0", Method: "job", Params: result}
// 	return st.enc.Encode(&message)
// }

// func (st *StratumSession) handleMessage(r *RpcServer, req *JSONRpcReq) error {
// 	if req.Id == nil {
// 		err := fmt.Errorf("Request ID Null")
// 		log.Error().Err(err)
// 		return err
// 	} else if req.Params == nil {
// 		err := fmt.Errorf("Request Params Empty")
// 		log.Error().Err(err)
// 		return err
// 	}

// 	// Handle RPC methods
// 	switch req.Method {
// 	case "login":
// 		r.registerSession(st)
// 		log.Debug().Msg(fmt.Sprintf("Login from %s", st.ip))
// 		miner := struct{
// 			uuid string `json:"login"`
// 		}{}
// 		err := json.Unmarshal(req.Params, &miner)
// 		if err != nil {
// 			return err
// 		}
// 		var loginJob := stratum.Login(miner.uuid)
// 		// Save data so we can recall for job
// 		st.minerId = miner.uuid
// 		return st.sendResult(req.Id, loginJob, nil)
// 	case "getjob":
// 		return st.sendResult(req.Id, nil, nil)
// 	case "submit":
// 		log.Debug().Msg(fmt.Sprintf("Submission from %s", st.ip))
// 		result := stratum.Submit(req.Params)
// 		return st.sendResult(req.Id, result, nil)
// 	case "keepalived":
// 		return st.sendResult(req.Id, &StatusReply{Status: "KEEPALIVED"}, nil)
// 	default:
// 		// Should actually mark as error true
// 		return st.sendResult(req.Id, &ErrorReply{Code: -1, Message: "Invalid method"}, nil)
// 	}
// }
