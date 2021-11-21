package server

import (
	"net"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/ybbus/jsonrpc"
)

type PoolServer struct {
	Client     jsonrpc.RPCClient
	SessionsMu sync.RWMutex
	Sessions   map[*StratumSession]struct{}
}

var once sync.Once
var singleInstance *PoolServer

func GetServerInstance() *PoolServer {
	if singleInstance == nil {
		once.Do(
			func() {
				singleInstance = New()
			})
	}
	return singleInstance
}

func New() *PoolServer {
	poolServer := &PoolServer{}
	poolServer.Sessions = make(map[*StratumSession]struct{})
	return poolServer
}

func (s *PoolServer) Listen(bindAddr string) {
	addr, err := net.ResolveTCPAddr("tcp", bindAddr)
	if err != nil {
		log.Fatal().Err(err)
		return
	}
	server, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal().Err(err)
		return
	}
	defer server.Close()

	log.Info().Msgf("Server listening on %s", bindAddr)
	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			// Continue attempting to receive requests
			log.Fatal().Err(err)
			continue
		}
		conn.SetKeepAlive(true)
		ip, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
		st := &StratumSession{conn: conn, ip: ip}
		s.registerSession(st)
		go st.handleSession()
	}
}

func (s *PoolServer) removeSession(session *StratumSession) {
	s.SessionsMu.Lock()
	defer s.SessionsMu.Unlock()
	delete(s.Sessions, session)
}
func (s *PoolServer) registerSession(session *StratumSession) {
	s.SessionsMu.Lock()
	defer s.SessionsMu.Unlock()
	s.Sessions[session] = struct{}{}
}
