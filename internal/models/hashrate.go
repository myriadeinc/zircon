package models

type MinerHashrate struct {
	MinerId         string `json:"minerId"`
	TotalDifficulty uint64 `json:"totalDifficulty"`
}
