package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ybbus/jsonrpc"
)

const (
	// 1024 bytes * 10 = 10kb max size
	MaxReqSize = 10 * 1024
	// For now, 10k at a time
	MaxConns = 10000
)

type RpcServer struct {
	Client     jsonrpc.RPCClient
	sessionsMu sync.RWMutex
	sessions   map[*StratumSession]struct{}
}

type StratumSession struct {
	sync.Mutex
	conn    *net.TCPConn
	enc     *json.Encoder
	ip      string
	minerId string
}

func New() *RpcServer {
	rpcServer := &RpcServer{}
	rpcServer.sessions = make(map[*StratumSession]struct{})
	rpcServer.NewClient()
	return rpcServer
}

func (r *RpcServer) NewClient() {
	r.Client = jsonrpc.NewClient("http://localhost:22345")
}

type TestResp struct {
	Status string `json:"status"`
}

type Request struct {
	Wallet_Address string      `json:"wallet_address"`
	Reserve_Size   int         `json:"reserve_size"`
	Other          interface{} `json:"more"`
}

type JSONRpcReq struct {
	Id     *json.RawMessage `json:"id"`
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
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

func (r *RpcServer) ListenHTTP() {
	var port string = ":4990"
	http.HandleFunc("/jobfeed", r.Jobfeed)
	log.Info().Msgf("Listening on %s%s:%s", os.Getenv("service__local__name"), port, "/jobfeed")
	http.ListenAndServe(port, nil)

}

func (r *RpcServer) Listen(bindAddr string) {
	// bindAddr := fmt.Sprintf("%s:%d", e.config.Host, e.config.Port)
	// Init TCP server
	// var port int = 8222
	// bindAddr := fmt.Sprintf("%s:%d", "0.0.0.0", port)
	addr, err := net.ResolveTCPAddr("tcp", bindAddr)
	if err != nil {
		log.Fatal().Err(err)
	}
	server, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer server.Close()

	log.Info().Msgf("Server listening on %s", bindAddr)
	// By passing the channel, we are handling (MaxConns) connections at a time
	accept := make(chan int, MaxConns)
	n := 0
	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			log.Fatal().Err(err)
			continue
		}
		conn.SetKeepAlive(true)
		ip, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
		st := &StratumSession{conn: conn, ip: ip, enc: json.NewEncoder(conn)}
		n++
		accept <- n
		go func() {
			r.handleClient(st)
			<-accept
		}()
	}
}

func (r *RpcServer) handleClient(st *StratumSession) {
	connbuff := bufio.NewReaderSize(st.conn, MaxReqSize)
	//stet max connection timeout
	duration := time.Second * time.Duration(360)
	st.conn.SetDeadline(time.Now().Add(duration))

	for {
		data, isPrefix, err := connbuff.ReadLine()
		if isPrefix {
			log.Error().Msgf("Socket flood detected from %s", st.ip)
			break
		} else if err == io.EOF {
			log.Error().Msgf("Client disconnected %s", st.ip)
			break
		} else if err != nil {
			log.Error().Msg(fmt.Sprintf("Error reading:", err))
			break
		}

		var req JSONRpcReq
		err = json.Unmarshal(data, &req)

		if err != nil {
			break
		}
		st.conn.SetDeadline(time.Now().Add(duration))

		err = st.handleMessage(r, &req)
		if err != nil {
			break
		}

	}
	r.removeSession(st)
	st.conn.Close()
}

func (r *RpcServer) Jobfeed(w http.ResponseWriter, req *http.Request) {
	r.broadcastNewJobs()
	fmt.Fprintf(w, "OK\n")
}

func (r *RpcServer) broadcastNewJobs() {
	r.sessionsMu.RLock()
	defer r.sessionsMu.RUnlock()
	count := len(r.sessions)
	log.Debug().Msgf("Broadcasting new jobs to %d miners", count)
	// Why 1024 * 16? 16 bytes?
	bcast := make(chan int, 1024*16)
	n := 0

	for m := range r.sessions {
		n++
		bcast <- n
		go func(s *StratumSession) {
			// job := r.callEmerald("proxyjob", cs.LoginData)
			// msg := handleEmeraldJob(job, false)

			// err := cs.sendJob(msg)
			err := s.sendJob(nil)
			<-bcast
			if err != nil {
				log.Error().Msgf("Job transmit error to %s: %v", s.ip, err)
				r.removeSession(s)
			}
			duration := time.Second * time.Duration(360)
			s.conn.SetDeadline(time.Now().Add(duration))
		}(m)
	}

}
