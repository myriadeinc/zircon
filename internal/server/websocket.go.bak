package server

import (
	b64 "encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/myriadeinc/zircon/internal/stratum"
	"github.com/rs/zerolog/log"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func handleMessageB64(message []byte) (bool, string, error) {
	decoded, _ := b64.StdEncoding.DecodeString(string(message))
	var request JSONRpcReq
	var minerId string
	jsonErr := json.Unmarshal(decoded, &request)
	if jsonErr != nil {
		return false, "", jsonErr
	}
	if request.Method == "login" {
		minerId = request.ParseMinerId()
	}
	needNewJob, response, stratumErr := request.GetStratumResponse(minerId)
	if stratumErr != nil {
		return false, "", stratumErr
	}
	encodedPayload := b64.StdEncoding.EncodeToString(response)
	return needNewJob, encodedPayload, nil
}

func Echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Error().Msgf("read:", err)
			break
		}
		log.Info().Msgf("receive: %s", message)
		needNewJob, payload, err := handleMessageB64(message)
		if err != nil {
			log.Error().Msgf("write:", err)
			break
		}
		if needNewJob {
			job := stratum.GetDummyJobWrapper()
			bytes, _ := json.Marshal(job)
			encodedPayload := b64.StdEncoding.EncodeToString(bytes)
			err = c.WriteMessage(mt, []byte(encodedPayload))
			if err != nil {
				log.Error().Msgf("write:", err)
				break
			}
		} else {
			err = c.WriteMessage(mt, []byte(payload))
			if err != nil {
				log.Error().Msgf("write:", err)
				break
			}
		}

	}
}
func MockStratum(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Msgf("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Error().Msgf("read:", err)
			break
		}
		log.Info().Msgf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Error().Msgf("write:", err)
			break
		}
	}
}
func test() {
	http.HandleFunc("/debug", Echo)
	http.ListenAndServe("localhost:8080", nil)
}
