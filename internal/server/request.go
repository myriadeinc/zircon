package server

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/myriadeinc/zircon/internal/stratum"
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

func (r *JSONRpcReq) GetStratumResponse(minerId string) (bool, []byte, error) {
	needNewJob := false
	var jsonBody []byte
	var jsonErr error

	switch r.Method {
	case "login":
		// log.Debug().Msgf("Login call for %s", minerId)
		job := stratum.GetDummyResponse(r.Id)
		jsonBody, jsonErr = json.Marshal(job)
	case "submit":
		// log.Debug().Msgf("Submit call for %s", minerId)
		ok := stratum.GetDummyOk(r.Id)
		jsonBody, jsonErr = json.Marshal(ok)
		needNewJob = true
	default:
		unkownMethod := struct {
			Message string `json:"message"`
		}{Message: "unknownMethod"}
		jsonBody, jsonErr = json.Marshal(unkownMethod)
	}

	return needNewJob, jsonBody, jsonErr
}
