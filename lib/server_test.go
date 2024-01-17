package server

import "testing"

func TestRun(t *testing.T) {
	got := true
	if !got {
		t.Errorf("Failed because some reason")
	}
}