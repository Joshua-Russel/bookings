package main

import "testing"

func TestRun(t *testing.T) {
	db, err := run()
	if err != nil {
		t.Error("failed test")
	}

}
