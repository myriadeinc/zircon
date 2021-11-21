package server

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"sync"

	"github.com/myriadeinc/zircon/internal/stratum"
	"github.com/rs/zerolog/log"
)

type StratumSession struct {
	sync.Mutex
	conn    *net.TCPConn
	ip      string
	minerId string
}

func (s *StratumSession) handleSession() {
	defer s.conn.Close()
	defer GetServerInstance().removeSession(s)
	clientReader := bufio.NewReader(s.conn)
	for {
		// Could be long buffer, but at least we handle errors
		clientRequest, err := clientReader.ReadBytes('\n')
		switch err {
		case nil:
			requestErr := s.handleRequest(clientRequest)
			if requestErr != nil {
				log.Error().Msg("Disconnecting due to request process error")
				return
			}
		case io.EOF:
			log.Debug().Msg("Client disconnected")
			return
		default:
			log.Fatal().Msg("Error connectiuon")
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
	needNewJob, response, stratumErr := request.GetStratumResponse()
	if stratumErr != nil {
		return stratumErr
	}

	response = append(response, byte('\n'))
	_, writeErr := s.conn.Write(response)

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
