package stratum

import (
	. "gopkg.in/check.v1"
)

type StratumServiceSuite struct{}

var _ = Suite(&StratumServiceSuite{})

func (s *StratumServiceSuite) TestSaveAndFetch(c *C) {

	// h := NewStratumRPCService()
	// job := h.HandleCall()
	// // g := NewRedisService()
	// // g.DebugSave()
	// // t := g.GetBlockTemplate("24")
	// // fmt.Println(t)

	// c.Assert(job.Result, Equals, map[string]string{})
	// c.Assert("s", Equals, "l")
}
