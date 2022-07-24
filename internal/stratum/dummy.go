package stratum

import "encoding/json"

type DummyStratumService struct {
}

func GetDummyStratumService() StratumService {
	dummyService := &DummyStratumService{}
	return dummyService
}

func (d *DummyStratumService) HandleLogin(id *json.RawMessage, minerId string) (*LoginResponse, error) {

	return &LoginResponse{
		Id:      id,
		Version: "2.0",
		Result: map[string]interface{}{
			"id": "12345",
			"job": map[string]interface{}{
				"height":    "2498341",
				"blob":      "0e0e9bd8ea8c068d8a87fac40b47fd013e892e7b9d245191ded7f38796cb3559e7546186e894860000000034263929525d40577ccd6729494a95b4d9c01c3478e1120ba9b44dd2471013801b",
				"job_id":    "2498341",
				"target":    "b2df0000",
				"algo":      "rx/0",
				"seed_hash": "d02e1b1704b67497736d3d3bb25423c3fab9e472536f1353cf57f9014068ffd8",
				"status":    "OK",
			},
		},
	}, nil
}

func (d *DummyStratumService) HandleSubmit(id *json.RawMessage, result *json.RawMessage) (*SubmitResponse, error) {
	return &SubmitResponse{
		Id:      id,
		Version: "2.0",
		Result: map[string]string{
			"status": "OK",
		},
	}, nil
}
func (d *DummyStratumService) HandleNewJob(minerId string) (*JobResponse, error) {
	return &JobResponse{
		Version: "2.0",
		Method:  "job",
		Params: Job{
			"Height":   "3498341",
			"Blob":     "0e0e9bd8ea8c068d8a87fac40b47fd013e892e7b9d245191ded7f38796cb3559e7546186e894860000000034263929525d40577ccd6729494a95b4d9c01c3478e1120ba9b44dd2471013801b",
			"JobId":    RandomInt(),
			"Target":   "b2df0000",
			"Algo":     "rx/0",
			"SeedHash": "d02e1b1704b67497736d3d3bb25423c3fab9e472536f1353cf57f9014068ffd8",
			"Status":   "OK",
		},
	}, nil
}
