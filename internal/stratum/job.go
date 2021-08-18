package stratum

import (
	"context"
)




func (job *JobData) updateTarget() JobData {
	// Parse as int64 in base 10
	target, _ := strconv.ParseInt(job.Target, 10, 64)
	job.Target = GetTargetHex(target)
	return *job
}


func CreateJob(ctx context.Context) {
	bt := GetBlockTemplate()
// create nonce
 newJob := Job{
	minerId: minerId,
	job_id: ,// random 64 chars as hex
	extraNonce: extraNonce,
	height: blockTemplate.height,
	seed_hash: blockTemplate.seed_hash,
	blob: blockTemplate.templateBlob,
	globalDiff: 

	difficulty: 
  };
	cache.Set()
}
func FetchJob(ctx context.Context) {}
