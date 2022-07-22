package server

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/rs/zerolog/log"
)

type JSONRpcReq struct {
	Id     *json.RawMessage `json:"id"`
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
}

// Set a default minerId so that we can collect metrics down the line
const funnelMinerId = "00001111-1111-4222-8333-abc123456789"

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func (r *JSONRpcReq) ParseMinerId() string {
	loginData := make(map[string]interface{})
	err := json.Unmarshal(*r.Params, &loginData)
	if _, exists := loginData["login"]; !exists || err != nil {
		log.Error().Msg("Could not parse minerId from login")
		return funnelMinerId
	}
	minerId := fmt.Sprintf("%v", loginData["login"])

	if !isValidUUID(minerId) {
		return funnelMinerId
	}

	return minerId

}
