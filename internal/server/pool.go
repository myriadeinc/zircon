package server

import (
	"net"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/ybbus/jsonrpc"
)

type PoolServer struct {
	Client jsonrpc.RPCClient
}

type SessionHandler struct {
	SessionsMu sync.RWMutex
	Sessions   map[*StratumSession]struct{}
}

var once sync.Once
var sessionHandler *SessionHandler

func GetSessionHandler() *SessionHandler {
	if sessionHandler == nil {
		once.Do(
			func() {
				sessionMap := make(map[*StratumSession]struct{})
				sessionHandler = &SessionHandler{
					Sessions: sessionMap,
				}
			})
	}
	return sessionHandler
}

func New() *PoolServer {
	poolServer := &PoolServer{}
	_ = GetSessionHandler()
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

	log.Info().Msgf("Server started! listening on %s", bindAddr)
	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			log.Fatal().Err(err)
			return
		}
		conn.SetKeepAlive(true)
		ip, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
		st := NewSession(ip, conn)
		sessions := GetSessionHandler()
		sessions.registerSession(st)
		go st.handleSession()
	}
}

func (s *SessionHandler) removeSession(session *StratumSession) {
	s.SessionsMu.Lock()
	defer s.SessionsMu.Unlock()
	delete(s.Sessions, session)
}
func (s *SessionHandler) registerSession(session *StratumSession) {
	s.SessionsMu.Lock()
	defer s.SessionsMu.Unlock()
	s.Sessions[session] = struct{}{}
}
