package server

import (
	"encoding/json"

	"github.com/myriadeinc/zircon/internal/stratum"
	"github.com/rs/zerolog/log"
)

type JSONRpcReq struct {
	Id     *json.RawMessage `json:"id"`
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
}

// Set a default minerId so that we can collect metrics down the line
const funnelMinerId = "699f06ec-9756-4d42-a470-7bf5f104f0e2"

func (r *JSONRpcReq) ParseMinerId() string {
	loginData := make(map[string]string)
	err := json.Unmarshal(*r.Params, &loginData)
	if _, exists := loginData["login"]; !exists || err != nil {
		log.Error().Msg("Could not parse minerId")
		return funnelMinerId
	}
	return loginData["login"]

}

func (r *JSONRpcReq) GetStratumResponse() (bool, []byte, error) {
	needNewJob := false
	var jsonBody []byte
	var jsonErr error

	switch r.Method {
	case "login":
		log.Info().Msg("Login detected")
		job := stratum.GetDummyResponse(r.Id)
		jsonBody, jsonErr = json.Marshal(job)
	case "submit":
		log.Info().Msg("Submit detected")
		ok := stratum.GetDummyOk
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
