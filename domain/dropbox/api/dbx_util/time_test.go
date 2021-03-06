package dbx_util

import (
	"testing"
	"time"
)

func TestRebaseTimeForAPI(t *testing.T) {
	jst, err := time.LoadLocation("Japan")
	if err != nil {
		t.Error(err)
	}
	nowUtc := time.Now()
	nowJst := nowUtc.In(jst)
	nowRoundedUtc := nowUtc.Round(1000 * time.Millisecond)

	if !RebaseTime(nowJst).Equal(nowRoundedUtc) {
		t.Error("Invalid state")
	}
}
