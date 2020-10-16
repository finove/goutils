package redisop

import (
	"fmt"
	"testing"

	"github.com/garyburd/redigo/redis"
)

func testLog() {
	var ll = NewDefaultLogger("debug")
	SetLogger(ll)
	warn("hello, warn %d,%d,%d\n", levelWarn, levelInfo, levelDebug)
	info("hello, info\n")
	debug("hello, debug %s\n", "ddd")
}

func testRedis() {
	var r = &RedisPoolX{}
	r.Connect("default")
	r.Rdo("SET", "tmpa", "HELLO WORLD")
	vv, err := redis.String(r.Rdo("GET", "tmpa"))
	fmt.Printf("err %v, value %v\n", err, vv)
	fmt.Printf("redis notify config %s\n", r.GetNotifyConfig())
}

func testSub() {
	var p = &RedisPoolX{}
	p.Connect("default")
	p.Subscribe("__keyspace@0__:aa*", true)
	go p.SubscribeHandle(defaultPubSubHandle)
	select {}
}

func TestAll(t *testing.T) {
	testLog()
	// testRedis()
	testSub()
}
