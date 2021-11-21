package server

import (
	"encoding/json"
)

type TestResp struct {
	Status string `json:"status"`
}

type Request struct {
	Wallet_Address string      `json:"wallet_address"`
	Reserve_Size   int         `json:"reserve_size"`
	Other          interface{} `json:"more"`
}

type JSONRpcResp struct {
	Id      *json.RawMessage `json:"id"`
	Version string           `json:"jsonrpc"`
	Result  interface{}      `json:"result"`
	Error   interface{}      `json:"error"`
}

// Push Job to XMRig
type JSONRpcPushMessage struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

// func (r *RpcServer) ListenHTTP() {
// 	var port string = ":4990"
// 	http.HandleFunc("/jobfeed", r.Jobfeed)
// 	log.Info().Msgf("Listening on %s%s:%s", os.Getenv("service__local__name"), port, "/jobfeed")
// 	http.ListenAndServe(port, nil)
// }

// func (r *RpcServer) handleClient(st *StratumSession) {
// 	connbuff := bufio.NewReaderSize(st.conn, MaxReqSize)
// 	//stet max connection timeout
// 	duration := time.Second * time.Duration(360)
// 	st.conn.SetDeadline(time.Now().Add(duration))

// 	for {
// 		data, isPrefix, err := connbuff.ReadLine()
// 		if isPrefix {
// 			log.Error().Msgf("Socket flood detected from %s", st.ip)
// 			break
// 		} else if err == io.EOF {
// 			log.Error().Msgf("Client disconnected %s", st.ip)
// 			break
// 		} else if err != nil {
// 			log.Error().Msg(fmt.Sprintf("Error reading:", err))
// 			break
// 		}

// 		var req JSONRpcReq
// 		err = json.Unmarshal(data, &req)

// 		if err != nil {
// 			break
// 		}
// 		st.conn.SetDeadline(time.Now().Add(duration))

// 		err = st.handleMessage(r, &req)
// 		if err != nil {
// 			break
// 		}

// 	}
// 	r.removeSession(st)
// 	st.conn.Close()
// }

// func (r *RpcServer) Jobfeed(w http.ResponseWriter, req *http.Request) {
// 	r.broadcastNewJobs()
// 	fmt.Fprintf(w, "OK\n")
// }

// func (r *RpcServer) broadcastNewJobs() {
// 	r.sessionsMu.RLock()
// 	defer r.sessionsMu.RUnlock()
// 	count := len(r.sessions)
// 	log.Debug().Msgf("Broadcasting new jobs to %d miners", count)
// 	// Why 1024 * 16? 16 bytes?
// 	bcast := make(chan int, 1024*16)
// 	n := 0

// 	for m := range r.sessions {
// 		n++
// 		bcast <- n
// 		go func(s *StratumSession) {
// 			// job := r.callEmerald("proxyjob", cs.LoginData)
// 			// msg := handleEmeraldJob(job, false)

// 			// err := cs.sendJob(msg)
// 			err := s.sendJob(nil)
// 			<-bcast
// 			if err != nil {
// 				log.Error().Msgf("Job transmit error to %s: %v", s.ip, err)
// 				r.removeSession(s)
// 			}
// 			duration := time.Second * time.Duration(360)
// 			s.conn.SetDeadline(time.Now().Add(duration))
// 		}(m)
// 	}

// }
