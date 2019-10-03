package rest

import (
	"testing"
	"time"
)

func TestRateLimit(t *testing.T) {
	r := RateLimit{
		Limit:     10,
		Remaining: 10,
		Period:    10,
	}
	if r.PercentageLeft() != 100 {
		t.Error("PercentLeft != 100")
	}
	if r.WaitTime() != time.Second {
		t.Error("WaitTime is wrong duration ", r.WaitTime())
	}
	if r.WaitTimeRemaining() != (time.Duration(1) * time.Second) {
		t.Error("WaitTimeRemaining is wrong duration ", r.WaitTimeRemaining())
	}

	r.Remaining = 5
	if r.PercentageLeft() != 50 {
		t.Error("PercentLeft != 50")
	}
	if r.WaitTime() != time.Second {
		t.Error("WaitTime is wrong duration ", r.WaitTime())
	}
	if r.WaitTimeRemaining() != (time.Duration(2) * time.Second) {
		t.Error("WaitTimeRemaining is wrong duration ", r.WaitTimeRemaining())
	}

	r.Remaining = 0
	if r.PercentageLeft() != 0 {
		t.Error("PercentLeft != 0")
	}
	if r.WaitTime() != time.Second {
		t.Error("WaitTime is wrong duration ", r.WaitTime())
	}
	if r.WaitTimeRemaining() != (time.Duration(10) * time.Second) {
		t.Error("WaitTimeRemaining is wrong duration ", r.WaitTimeRemaining())
	}
}
