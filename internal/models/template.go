package models

import (
	"encoding/json"
	"fmt"
)

// Instead of map[string]string because the redis library only has support to serialize a struct
type StrictTemplate struct {
	BlockTemplateBlob string `json:"blocktemplate_blob"`
	Difficulty        string `json:"difficulty"`
	SeedHash          string `json:"seed_hash"`
	Height            string `json:"height"`
}

func (s StrictTemplate) IsValid() bool {
	return len(s.BlockTemplateBlob) > 0 && len(s.SeedHash) > 0 && len(s.Difficulty) > 0 && len(s.Height) > 0
}

func (s StrictTemplate) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s StrictTemplate) ToLoginRequest(minerId string) map[string]string {
	return map[string]string{
		"miner":            minerId,
		"templateBlob":     s.BlockTemplateBlob,
		"templateHeight":   s.Height,
		"templateDiff":     s.Difficulty,
		"templateSeedhash": s.SeedHash,
	}
}

func NewStrictTemplate(template map[string]interface{}) (*StrictTemplate, error) {
	if len(template) == 0 {
		return nil, fmt.Errorf("template is empty")
	}

	stemplate := StrictTemplate{
		BlockTemplateBlob: fmt.Sprintf("%v", template["blocktemplate_blob"]),
		SeedHash:          fmt.Sprintf("%v", template["seed_hash"]),
		// Horrible yes, but desired as they are parsed as floats but should be uint
		Difficulty: fmt.Sprintf("%d", uint64(template["difficulty"].(float64))),
		Height:     fmt.Sprintf("%d", uint64(template["height"].(float64))),
	}
	if !stemplate.IsValid() {
		return nil, fmt.Errorf("expected full template, received %v", template)
	}

	return &stemplate, nil
}
