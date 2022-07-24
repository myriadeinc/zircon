package stratum

import (
	. "gopkg.in/check.v1"
)

type UtilSuite struct{}

var _ = Suite(&UtilSuite{})

func (s *UtilSuite) TestGenerateHex500(c *C) {

	targetHex := convertDifficultyToHex("500")
	expectedHex := "6e128300"

	c.Assert(targetHex, Equals, expectedHex)
}

func (s *UtilSuite) TestGenerateHex15k(c *C) {

	targetHex := convertDifficultyToHex("15000")
	expectedHex := "7b5e0400"

	c.Assert(targetHex, Equals, expectedHex)
}
