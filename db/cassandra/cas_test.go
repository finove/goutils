package cassandra

import (
	"testing"
)

func TestAll(t *testing.T) {
	var err error
	var ag agent
	err = ag.Connect("brig_test", "127.0.0.1")
	if err != nil {
		t.Logf("connect fail %v", err)
		t.Fail()
		return
	}
}
