package server

import (
	"net"

	"github.com/myriadeinc/zircon/internal/stratum"

	"github.com/rs/zerolog/log"
	"github.com/ybbus/jsonrpc"
)

type PoolServer struct {
	Client  jsonrpc.RPCClient
	Stratum stratum.StratumService
}

func New() *PoolServer {
	poolServer := &PoolServer{
		Stratum: stratum.NewStratumRPCService(),
	}
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

	log.Info().Msgf("xmrig compatible server started! listening on %s", bindAddr)
	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			log.Fatal().Err(err)
			return
		}
		// Fire and forget sessions
		conn.SetKeepAlive(true)
		ip, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
		st := NewSession(ip, conn, s.Stratum)
		sessions := GetSessionHandler()
		sessions.addSession(st)
		go st.handleSession()
	}
}
