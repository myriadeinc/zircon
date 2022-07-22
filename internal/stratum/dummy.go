package stratum

import "encoding/json"

type DummyStratumService struct {
}

func GetDummyStratumService() StratumService {
	dummyService := &DummyStratumService{}
	return dummyService
}

func (d *DummyStratumService) HandleLogin(id *json.RawMessage) LoginResponse {

	job := GetDummyJobWrapper()
	return LoginResponse{
		Id:      id,
		Version: "2.0",
		Result:  job,
	}
}

func (d *DummyStratumService) HandleSubmit(id *json.RawMessage) (SubmitResponse, bool) {
	return SubmitResponse{
		Id:      id,
		Version: "2.0",
		Result: map[string]string{
			"status": "OK",
		},
	}, true
}
func (d *DummyStratumService) HandleNewJob() JobResponse {
	return JobResponse{
		Version: "2.0",
		Method:  "job",
		Params: Job{
			Height:   2498341,
			Blob:     "0e0e9bd8ea8c068d8a87fac40b47fd013e892e7b9d245191ded7f38796cb3559e7546186e894860000000034263929525d40577ccd6729494a95b4d9c01c3478e1120ba9b44dd2471013801b",
			JobId:    RandomInt(),
			Target:   "b2df0000",
			Algo:     "rx/0",
			SeedHash: "d02e1b1704b67497736d3d3bb25423c3fab9e472536f1353cf57f9014068ffd8",
			Status:   "OK",
		},
	}
}
