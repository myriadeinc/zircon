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

func (st *StratumSession) sendResult(id *json.RawMessage, result interface{}, err error) error {
	if err != nil {
		message := JSONRpcResp{Id: id, Version: "2.0", Error: err, Result: nil}
		return st.enc.Encode(&message)
	}
	st.Lock()
	defer st.Unlock()
	message := JSONRpcResp{Id: id, Version: "2.0", Error: nil, Result: result}
	return st.enc.Encode(&message)
}

// Push notification
func (st *StratumSession) sendJob(result interface{}) error {
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
		miner := struct{
			uuid string `json:"login"`
		}{}
		err := json.Unmarshal(req.Params, &miner)
		if err != nil {
			return err
		}
		var loginJob := stratum.Login(miner.uuid)
		// Save data so we can recall for job
		st.minerId = miner.uuid
		return st.sendResult(req.Id, loginJob, nil)
	case "getjob":
		return st.sendResult(req.Id, nil, nil)
	case "submit":
		log.Debug().Msg(fmt.Sprintf("Submission from %s", st.ip))
		result := stratum.Submit(req.Params)
		return st.sendResult(req.Id, result, nil)
	case "keepalived":
		return st.sendResult(req.Id, &StatusReply{Status: "KEEPALIVED"}, nil)
	default:
		// Should actually mark as error true
		return st.sendResult(req.Id, &ErrorReply{Code: -1, Message: "Invalid method"}, nil)
	}
}
