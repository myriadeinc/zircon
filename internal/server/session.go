package server

import (
	"encoding/json"
	"net"
	"strings"

	"github.com/myriadeinc/zircon/internal/cache"
	"github.com/myriadeinc/zircon/internal/stratum"
)

type StratumSession struct {
	conn    *net.TCPConn
	ip      string
	minerId string
}

func NewSession(ip string, connection *net.TCPConn, service stratum.StratumService, cache cache.CacheService) *StratumSession {
	return &StratumSession{
		ip:   ip,
		conn: connection,
	}
}

func isEmptyRequest(rawRequest []byte) bool {
	s := strings.TrimSpace(string(rawRequest))
	return len(s) == 0
}

func genericErrorResponse(id *json.RawMessage) []byte {
	response := map[string]interface{}{
		"id":      id,
		"jsonrpc": "2.0",
		"error": map[string]interface{}{
			"code":    -1,
			"message": "Internal server error",
		},
	}

	bytes, _ := json.Marshal(response)
	return bytes

}
