package rawserver

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ybbus/jsonrpc"
)

const (
	// 10 kB ?
	MaxReqSize = 10 * 1024
)

type RpcServer struct {
	Client     jsonrpc.RPCClient
	sessionsMu sync.RWMutex
	sessions   map[*Session]struct{}
}

type Session struct {
	sync.Mutex
	conn      *net.TCPConn
	enc       *json.Encoder
	ip        string
	LoginData *json.RawMessage
}

type Endpoint struct {
	id uint32
}

func New() *RpcServer {
	rpcServer := &RpcServer{}
	rpcServer.sessions = make(map[*Session]struct{})
	rpcServer.NewClient()

	return rpcServer
}

func NewEndpoint() *Endpoint {
	// Later we will map id to difficulty
	e := &Endpoint{id: 5678}
	return e
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

func (r *RpcServer) callEmerald(method string, rawRequest *json.RawMessage) *json.RawMessage {

	rpc := jsonrpc.NewClient("http://emerald:22345")

	var proxyResp *json.RawMessage

	response, _ := rpc.Call(method, rawRequest)
	err := response.GetObject(&proxyResp)
	if proxyResp == nil {
		log.Error().Msg("Emerald returned blank")
		return nil

	} else if err != nil {
		log.Error().Msg("Error contacting emerald service")
		log.Error().Err(err)
		return nil
	}
	return proxyResp

}

func (r *RpcServer) ListenHTTP() {
	var port string = ":4990"
	http.HandleFunc("/jobfeed", r.Jobfeed)
	log.Info().Msgf("Listening on %s%s:%s", os.Getenv("service__local__name"), port, "/jobfeed")
	http.ListenAndServe(port, nil)

}

// this does stuff
// you know
func (r *RpcServer) Listen() {

	quit := make(chan bool)
	// Wrapped to control RPC Server, also for future versions with multiple ports if necessary
	go func() {
		e := NewEndpoint()
		e.Listen(r)
	}()
	<-quit
}

func (e *Endpoint) Listen(r *RpcServer) {
	// bindAddr := fmt.Sprintf("%s:%d", e.config.Host, e.config.Port)
	var port int = 8222
	bindAddr := fmt.Sprintf("%s:%d", "0.0.0.0", port)
	addr, err := net.ResolveTCPAddr("tcp", bindAddr)
	if err != nil {
		log.Fatal().Err(err)
	}
	server, err := net.ListenTCP("tcp", addr)

	if err != nil {
		log.Fatal().Err(err)
	}
	defer server.Close()

	// msg := fmt.Sprintf()
	log.Info().Msgf("Server listening on %s", bindAddr)
	maxConns := 1000
	accept := make(chan int, maxConns)
	n := 0

	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			continue
		}
		conn.SetKeepAlive(true)
		ip, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
		cs := &Session{conn: conn, ip: ip, enc: json.NewEncoder(conn)}
		n++

		accept <- n
		go func() {
			r.handleClient(cs)
			<-accept
		}()
	}
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
		go func(cs *Session) {
			job := r.callEmerald("proxyjob", cs.LoginData)
			msg := handleEmeraldJob(job, false)

			err := cs.sendJob(msg)
			<-bcast
			if err != nil {
				log.Error().Msgf("Job transmit error to %s: %v", cs.ip, err)
				r.removeSession(cs)
			}
			duration := time.Second * time.Duration(360)
			cs.conn.SetDeadline(time.Now().Add(duration))
		}(m)
	}

}
