package server

import (
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/myriadeinc/zircon/internal/cache"
	"github.com/myriadeinc/zircon/internal/nodeapi"
	"github.com/myriadeinc/zircon/internal/stratum"
	"github.com/spf13/viper"

	"github.com/rs/zerolog/log"
)

type PoolServer struct {
	NodeClient  nodeapi.NodeApi
	Stratum     stratum.StratumService
	Cache       cache.CacheService
	SessionLock sync.RWMutex
	Sessions    map[*StratumSession]struct{}
	blockHeight uint64
}

func New() *PoolServer {
	nodes := strings.Split(viper.GetString("MONERO_NODES"), ",")

	poolServer := &PoolServer{
		NodeClient:  nodeapi.NewNodeClient(nodes),
		Stratum:     stratum.NewStratumRPCService(),
		Cache:       cache.NewClient(),
		Sessions:    make(map[*StratumSession]struct{}),
		blockHeight: 0,
	}
	return poolServer
}

func (s *PoolServer) Start(bindAddr string) {
	addr, err := net.ResolveTCPAddr("tcp", bindAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("")
		return
	}
	server, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal().Err(err).Msg("")
		return
	}
	defer server.Close()
	s.initTemplate()

	log.Info().Msgf("xmrig compatible server starting! listening on %s", bindAddr)
	go s.pollForever() // for loop that calls another go routine to broadcast new jobs
	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			log.Fatal().Err(err).Msg("could not start server")
			return
		}
		// Fire and forget sessions
		conn.SetKeepAlive(true)
		ip, _, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			log.Fatal().Err(err).Msg("could not activate connection")
			return
		}
		log.Debug().Str("ip", ip).Msg("new connection")
		session := NewSession(ip, conn, s.Stratum, s.Cache)
		s.addSession(session)
		go s.handlePoolSession(session)
	}

}

func (s *PoolServer) initTemplate() {

	template, templateErr := s.NodeClient.GetValidBlockTemplate()
	if templateErr != nil {
		log.Fatal().Err(templateErr).Msg("bad response from node")
		return
	}

	height, convErr := strconv.ParseUint(template.Height, 10, 64)
	if templateErr == nil && convErr == nil {
		s.blockHeight = height
		err := s.Cache.SaveNewTemplate(*template)
		if err != nil {
			log.Error().Err(err).Msg("could not save template")
		}
	}
}
