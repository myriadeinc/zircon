package server

import (
	"encoding/json"
	"strconv"
)

type StatusReply struct {
	Status string `json:"status"`
}

type ErrorReply struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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
