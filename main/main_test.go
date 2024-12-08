package main

import "testing"

func TestRun(t *testing.T) {

	if _, _, err := run(&app); err != nil {
		t.Error("TestFailed")
	}
}
