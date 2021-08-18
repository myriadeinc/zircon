package server

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
)

func (r *RpcServer) removeSession(cs *StratumSession) {
	r.sessionsMu.Lock()
	defer r.sessionsMu.Unlock()
	delete(r.sessions, cs)
}
func (r *RpcServer) registerSession(cs *StratumSession) {
	r.sessionsMu.Lock()
	defer r.sessionsMu.Unlock()
	r.sessions[cs] = struct{}{}
}

func (st *StratumSession) sendResult(id *json.RawMessage, result *json.RawMessage, err error) error {
	if err != nil {
		errr := []byte(err.Error())
		message := JSONRpcResp{Id: id, Version: "2.0", Error: (*json.RawMessage)(&errr), Result: nil}
		return st.enc.Encode(&message)
	}
	st.Lock()
	defer st.Unlock()
	message := JSONRpcResp{Id: id, Version: "2.0", Error: nil, Result: result}
	return st.enc.Encode(&message)
}

// Push notification
func (st *StratumSession) sendJob(result *json.RawMessage) error {
	st.Lock()
	defer st.Unlock()
	message := JSONRpcPushMessage{Version: "2.0", Method: "job", Params: result}
	return st.enc.Encode(&message)
}

func (st *StratumSession) handleMessage(r *RpcServer, req *JSONRpcReq) error {
	if req.Id == nil {
		err := fmt.Errorf("Request ID Null")
		log.Error().Err(err)
		return err
	} else if req.Params == nil {
		err := fmt.Errorf("Request Params Empty")
		log.Error().Err(err)
		return err
	}

	// Handle RPC methods
	switch req.Method {
	case "login":
		r.registerSession(st)
		log.Debug().Msg(fmt.Sprintf("Login from %s", st.ip))
		miner := struct {
			Uuid string `json:"login"`
		}{}
		err := json.Unmarshal(*req.Params, &miner)
		if err != nil {
			return err
		}
		// Save data so we can recall for job
		st.minerId = miner.Uuid
		log.Info().Msg("login detected!")
		testresp := []byte(`{
		"id": "1be0b7b6-b15a-47be-a17d-46b2911cf7d0",
		"job": {
		  "blob": "0e0ee2839688068ec1e852279fe58b27163035d96a6feebd3517fedc0cd1e3d1e93dbb5405c1f800000000883c7e7175a14d4dfa2e776f7583818106e110dadf1079e2fe70adfdb67242d38201",
		  "job_id": "q7PLUPL25UV0z5Ij14IyMk8htXbj",
		  "target": "b88d0600",
		  "id":"1be0b7b6-b15a-47be-a17d-46b2911cf7d0",
		  "seed_hash": "e44c720334ac9f5a7cd72f38cac1e8bbcae66e63fb1ad1765823da87d287f959",
		  "algo": "rx/0"
		},
		"status": "OK"}`)

		return st.sendResult(req.Id, (*json.RawMessage)(&testresp), nil)
	case "getjob":
		return st.sendResult(req.Id, nil, nil)
	case "submit":
		log.Debug().Msg(fmt.Sprintf("Submission from %s", st.ip))
		// result := stratum.Submit(req.Params)
		return st.sendResult(req.Id, nil, nil)
	case "keepalived":
		return st.sendResult(req.Id, nil, nil)
	default:
		// Should actually mark as error true
		// return st.sendResult(req.Id, &ErrorReply{Code: -1, Message: "Invalid method"}, nil)
		return st.sendResult(req.Id, nil, nil)
	}
}
