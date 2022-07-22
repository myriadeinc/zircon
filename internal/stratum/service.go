package stratum

import "encoding/json"

type StratumService interface {
	HandleLogin(*json.RawMessage) LoginResponse

	HandleSubmit(*json.RawMessage) (SubmitResponse, bool)

	HandleNewJob() JobResponse
}

// var service StratumService
