package stratum

import (
	"encoding/json"
)

type Job struct {
	Blob     string `json:"blob"`
	JobID    string `json:"job_id"`
	Target   string `json:"target"`
	ID       string `json:"id"`
	Height   int    `json:"height"`
	SeedHash string `json:"seed_hash"`
	Algo     string `json:"algo"`
}

type LoginJob struct {
	Id     string   `json:"id"`
	Job    *JobData `json:"job"`
	Status string   `json:"status"`
}

type JobPushData struct {
	Jsonrpc string  `json:"jsonrpc"`
	Method  string  `json:"method"`
	Params  JobData `json:"params"`
}

func Login(minerId string) (json.RawMessage, error) {

	l := CreateJob(minerId)
	job, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(job), nil
}

func Submit(submitParams *json.RawMessage) *json.RawMessage {
	return nil
}

func GetJob(minerId string) *json.RawMessage {
	return nil
}
