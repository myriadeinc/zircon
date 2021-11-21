package stratum

import (
	"encoding/json"
)

type LoginRequest struct {
	MinerId string `json:"login"`
	Agent   string `json:"agent"`
}

type DummyOk struct {
	Id      *json.RawMessage `json:"id"`
	Version string           `json:"jsonrpc"`
	Result  OkMsg            `json:"result"`
}
type OkMsg struct {
	Status string `json:"status"`
}

type DummyResponse struct {
	Id      *json.RawMessage `json:"id"`
	Version string           `json:"jsonrpc"`
	Result  JobWrapper       `json:"result"`
}
type JobWrapper struct {
	Id  string `json:"id"`
	Job Job    `json:"job"`
}
type Job struct {
	Height   int64  `json:"height"`
	Blob     string `json:"blob"`
	JobId    string `json:"job_id"`
	Target   string `json:"target"`
	Algo     string `json:"algo"`
	SeedHash string `json:"seed_hash"`
	Status   string `json:"status"`
}

func GetDummyResponse(id *json.RawMessage) DummyResponse {
	return DummyResponse{
		Id:      id,
		Version: "2.0",
		Result: JobWrapper{
			Id: "1",
			Job: Job{
				Height:   2498341,
				Blob:     "0e0e9bd8ea8c068d8a87fac40b47fd013e892e7b9d245191ded7f38796cb3559e7546186e894860000000034263929525d40577ccd6729494a95b4d9c01c3478e1120ba9b44dd2471013801b",
				JobId:    "0",
				Target:   "b2df0000",
				Algo:     "rx/0",
				SeedHash: "d02e1b1704b67497736d3d3bb25423c3fab9e472536f1353cf57f9014068ffd8",
				Status:   "OK",
			},
		},
	}
}

func GetDummyJobWrapper() JobWrapper {
	return JobWrapper{
		Id: "1",
		Job: Job{
			Height:   2498341,
			Blob:     "0e0e9bd8ea8c068d8a87fac40b47fd013e892e7b9d245191ded7f38796cb3559e7546186e894860000000034263929525d40577ccd6729494a95b4d9c01c3478e1120ba9b44dd2471013801b",
			JobId:    "0",
			Target:   "b2df0000",
			Algo:     "rx/0",
			SeedHash: "d02e1b1704b67497736d3d3bb25423c3fab9e472536f1353cf57f9014068ffd8",
			Status:   "OK",
		},
	}
}

func GetDummyOk(id *json.RawMessage) DummyOk {
	return DummyOk{
		Id:      id,
		Version: "2.0",
		Result: OkMsg{
			Status: "OK",
		},
	}
}
