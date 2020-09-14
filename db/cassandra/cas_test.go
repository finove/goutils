package cassandra

import (
	"log"
	"testing"
)

func TestAll(t *testing.T) {
	var err error
	var ag Agent
	err = ag.Connect("brig_test", "127.0.0.1")
	if err != nil {
		t.Logf("connect fail %v", err)
		t.Fail()
		return
	}
	tbls := ag.Tables()
	log.Printf("%s have %d tables %v", ag.DefaultKeyspace, len(tbls), tbls)
}
