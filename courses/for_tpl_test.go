package courses

import (
	"testing"
)

func TestPastCurrentUpdcoming(t *testing.T) {
	// TODO various test input here

	total := len(byID)
	if total != len(Past())+len(Current())+len(Upcoming()) {
		t.Fatalf("sum id not correct")
	}
}
