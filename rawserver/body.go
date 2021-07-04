package rawserver

import (
	"encoding/json"
	"strconv"
)

type JSONRpcReq struct {
	Id     *json.RawMessage `json:"id"`
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
}

type JSONRpcResp struct {
	Id      *json.RawMessage `json:"id"`
	Version string           `json:"jsonrpc"`
	Result  interface{}      `json:"result"`
	Error   interface{}      `json:"error"`
}

// Push Job to XMRig
type JSONRpcPushMessage struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type LoginRequest struct {
	Login string `json:"login"`
	// Password string
}

type StatusReply struct {
	Status string `json:"status"`
}

type ErrorReply struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type JobData struct {
	Blob     string `json:"blob"`
	JobID    string `json:"job_id"`
	Target   string `json:"target"`
	ID       string `json:"id"`
	Height   int    `json:"height"`
	SeedHash string `json:"seed_hash"`
	Algo     string `json:"algo"`
}

type LoginReply struct {
	Id     string   `json:"id"`
	Job    *JobData `json:"job"`
	Status string   `json:"status"`
}

type LoginData struct {
	ID     string  `json:"id"`
	Job    JobData `json:"job"`
	Status string  `json:"status"`
}

type JobPushData struct {
	Jsonrpc string  `json:"jsonrpc"`
	Method  string  `json:"method"`
	Params  JobData `json:"params"`
}

func handleEmeraldJob(emerald *json.RawMessage, isLogin bool) interface{} {
	var jobData JobData
	if isLogin {
		var loginData LoginData
		json.Unmarshal(*emerald, &loginData)
		loginData.Job = loginData.Job.updateTarget()
		return loginData
	}
	json.Unmarshal(*emerald, &jobData)
	return jobData.updateTarget()
}

func (job *JobData) updateTarget() JobData {
	// Parse as int64 in base 10
	target, _ := strconv.ParseInt(job.Target, 10, 64)
	job.Target = GetTargetHex(target)
	return *job
}
