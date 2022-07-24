package stratum

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type StratumService interface {
	HandleLogin(*json.RawMessage, string) (*LoginResponse, error)

	HandleSubmit(*json.RawMessage, *json.RawMessage) (*SubmitResponse, error)

	HandleNewJob(string) (*JobResponse, error)
}

type StratumRPCService struct {
	patriciaClient *rpc.Client
}

func reconnect(waitTime int) *rpc.Client {
	time.Sleep(time.Duration(waitTime))
	client, err := rpc.Dial(viper.GetString("WS_RPC_URL"))
	if err != nil {
		log.Error().Err(err).Int("waitTime", waitTime).Msg("Reconnect triggered : could not contact websocket patricia, trying again")
		return reconnect(waitTime * 2)
	}
	return client
}

func NewStratumRPCService() StratumService {
	waitTime := 5000000000 // 5 seconds
	time.Sleep(time.Duration(waitTime))
	client, err := rpc.Dial(viper.GetString("WS_RPC_URL"))
	if err != nil {
		log.Error().Err(err).Msg("Could not contact websocket patricia")
		client = reconnect(waitTime)
	}
	log.Info().Msg("successfully connected to rpc websocket server")
	service := &StratumRPCService{
		patriciaClient: client,
	}

	return service
}

func (s *StratumRPCService) HandleLogin(id *json.RawMessage, minerId string) (*LoginResponse, error) {
	// We use maps instead of explicit structs to check for contents easily
	minerJob := map[string]string{}

	err := s.patriciaClient.Call(&minerJob, "newjob", map[string]string{"miner": minerId})
	if err != nil {
		log.Error().Err(err).Msg("could not contact patricia client")
		return nil, err
	}
	if _, ok := minerJob["target"]; !ok {
		log.Error().Msg("No target diff in miner job")
		return nil, errors.New("did not receive expected payload")
	}
	if len(minerJob) == 0 {
		return nil, errors.New("received empty payload")
	}

	minerJob["target"] = convertDifficultyToHex(minerJob["target"])

	log.Trace().Str("minerJob", fmt.Sprint(minerJob)).Msg("got minerJob")

	m := map[string]interface{}{
		"id":  minerJob["job_id"],
		"job": minerJob,
	}

	result := &LoginResponse{
		Id:      id,
		Version: "2.0",
		Result:  m,
	}
	return result, nil

}

// "params":{
// 	"id":"3220921a94dd7ebacc85bdbf508b23e6545c80fb81",
// 	"job_id":"3220921a94dd7ebacc85bdbf508b23e6545c80fb81",
// 	"nonce":"f1830100",
// 	"result":"c68384ce77a3f4b1ffacd7e94b42f7da827e46fd3e8dfba3caa5eacf6cca6a01"
// }
func (s *StratumRPCService) HandleSubmit(id *json.RawMessage, params *json.RawMessage) (*SubmitResponse, error) {

	job_params := map[string]string{}
	err := json.Unmarshal(*params, &job_params)
	if err != nil {
		log.Error().Err(err).Msg("Could not unmarshal submit job")
		return nil, err
	}
	response := map[string]bool{}
	err = s.patriciaClient.Call(&response, "submitjob", job_params)
	if err != nil {
		log.Error().Err(err).Msg("could not contact patricia client for submitjob")
		return nil, err
	}

	if accepted, ok := response["accepted"]; ok && accepted {
		submitOk := &SubmitResponse{
			Id:      id,
			Version: "2.0",
			Result: map[string]string{
				"status": "OK",
			},
		}

		return submitOk, nil
	}
	if len(response) == 0 {
		return nil, errors.New("received empty response")
	}

	return nil, errors.New("block not accepted")

}

func (s *StratumRPCService) HandleNewJob(minerId string) (*JobResponse, error) {
	minerJob := map[string]string{}

	err := s.patriciaClient.Call(&minerJob, "newjob", map[string]string{"miner": minerId})
	if err != nil {
		log.Error().Err(err).Msg("could not fetch new job")
		return nil, err
	}
	log.Trace().Str("minerJob", fmt.Sprint(minerJob)).Msg("Received new job to push from patricia")

	if len(minerJob) == 0 {
		return nil, errors.New("received empty payload")
	}

	minerJob["target"] = convertDifficultyToHex(minerJob["target"])

	newjob := &JobResponse{
		Version: "2.0",
		Method:  "job",
		Params:  minerJob,
	}

	return newjob, nil

}
