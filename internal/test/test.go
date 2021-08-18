package test

import (
	"encoding/json"
)

func GetRawJson() json.RawMessage {

	return []byte(`{"key":"value"}`)

}
