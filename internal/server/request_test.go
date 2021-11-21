package server

import (
	"encoding/json"
	"testing"
)

func TestAdd(t *testing.T) {

	newReq := &JSONRpcReq{
		Method: "login",
		Params: &json.RawMessage{},
	}
	newId := newReq.ParseMinerId()
	expected := funnelMinerId

	if newId != expected {
		t.Errorf("got %s, expected %s", newId, expected)
	}
}
