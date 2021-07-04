package rawserver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/rs/zerolog/log"
)

func (r *RpcServer) removeSession(cs *Session) {
	r.sessionsMu.Lock()
	defer r.sessionsMu.Unlock()
	delete(r.sessions, cs)
}
func (r *RpcServer) registerSession(cs *Session) {
	r.sessionsMu.Lock()
	defer r.sessionsMu.Unlock()
	r.sessions[cs] = struct{}{}
}
func (r *RpcServer) handleClient(cs *Session) {
	// Current buffer size 1024 * 10 10kB
	connbuff := bufio.NewReaderSize(cs.conn, MaxReqSize)
	duration := time.Second * time.Duration(360)
	cs.conn.SetDeadline(time.Now().Add(duration))

	for {
		data, isPrefix, err := connbuff.ReadLine()
		if isPrefix {
			log.Error().Msgf("Socket flood detected from %s", cs.ip)
			break
		} else if err == io.EOF {
			log.Error().Msgf("Client disconnected %s", cs.ip)
			break
		} else if err != nil {
			log.Error().Msg(fmt.Sprintf("Error reading:", err))
			break
		}

		// NOTICE: cpuminer-multi sends junk newlines, so we demand at least 1 byte for decode
		// NOTICE: Ns*CNMiner.exe will send malformed JSON on very low diff, not sure we should handle this
		if len(data) > 1 {
			var req JSONRpcReq
			err = json.Unmarshal(data, &req)

			if err != nil {
				break
			}
			cs.conn.SetDeadline(time.Now().Add(duration))
			err = cs.handleMessage(r, &req)
			if err != nil {
				break
			}
		}
	}
	r.removeSession(cs)
	cs.conn.Close()
}

func (cs *Session) sendResult(id *json.RawMessage, result interface{}, err error) error {
	if err != nil {
		message := JSONRpcResp{Id: id, Version: "2.0", Error: err, Result: nil}
		return cs.enc.Encode(&message)
	}
	cs.Lock()
	defer cs.Unlock()
	message := JSONRpcResp{Id: id, Version: "2.0", Error: nil, Result: result}
	return cs.enc.Encode(&message)
}

// Push notification
func (cs *Session) sendJob(result interface{}) error {
	cs.Lock()
	defer cs.Unlock()
	message := JSONRpcPushMessage{Version: "2.0", Method: "job", Params: result}
	return cs.enc.Encode(&message)
}

func (cs *Session) handleMessage(r *RpcServer, req *JSONRpcReq) error {
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
		r.registerSession(cs)
		log.Info().Msg(fmt.Sprintf("Login from %s", cs.ip))
		var fromEmerald *json.RawMessage = r.callEmerald("login", req.Params)
		msg := handleEmeraldJob(fromEmerald, true)
		// Save data so we can recall for job
		cs.LoginData = req.Params
		return cs.sendResult(req.Id, msg, nil)
	case "getjob":
		return cs.sendResult(req.Id, nil, nil)
	case "submit":
		log.Info().Msg(fmt.Sprintf("Submission from %s", cs.ip))
		result := r.callEmerald("submit", req.Params)
		return cs.sendResult(req.Id, result, nil)
	case "keepalived":
		return cs.sendResult(req.Id, &StatusReply{Status: "KEEPALIVED"}, nil)
	default:
		// Should actually mark as error true
		return cs.sendResult(req.Id, &ErrorReply{Code: -1, Message: "Invalid method"}, nil)
	}
}
