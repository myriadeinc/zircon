package rawserver

import (
	"testing"
)

func TestGetTargetHex(t *testing.T) {
	targetHex := GetTargetHex(500)
	expectedHex := "6e128300"
	if targetHex != expectedHex {
		t.Error("Invalid targetHex")
	}

	targetHex = GetTargetHex(15000)
	expectedHex = "7b5e0400"
	if targetHex != expectedHex {
		t.Error("Invalid targetHex")
	}
}
