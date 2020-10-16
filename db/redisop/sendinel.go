package redisop

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/garyburd/redigo/redis"
)

/*
SENTINEL get-master-addr-by-name mymaster
SENTINEL sentinels <master name> Show a list of sentinel instances for this master, and their state

*/

// RedisPoolX 支持哨兵
type RedisPoolX struct {
	RedisPool
	sentinelPool    *RedisPool
	sentinelServers []string
	sentinelMaster  string
	sentinelIdx     uint8
	sentinelMu      sync.Mutex
}

// Sentinel 连接哨兵服务器
func (p *RedisPoolX) Sentinel(servers, masterName string) (err error) {
	p.sentinelMu.Lock()
	p.sentinelMaster = masterName
	p.sentinelServers = make([]string, 0)
	for _, server := range strings.Split(servers, ",") {
		if server != "" {
			p.sentinelServers = append(p.sentinelServers, server)
		}
	}
	p.sentinelMu.Unlock()
	debug("sentinel server init %s %v", p.sentinelMaster, p.sentinelServers)
	if len(p.sentinelServers) == 0 {
		err = errors.New("no sentinel server configed")
		return
	}
	err = p.sentinelTryConnect()
	return
}

func (p *RedisPoolX) sentinelTryConnect() (err error) {
	var try int
	for try = 0; try < len(p.sentinelServers); try++ {
		err = p.sentinelToNext()
		if err == nil {
			break
		}
	}
	return
}

func (p *RedisPoolX) sentinelToNext() (err error) {
	var server, redisServer string
	p.sentinelMu.Lock()
	defer p.sentinelMu.Unlock()
	if p.sentinelPool != nil {
		p.sentinelPool.Close()
	}
	p.sentinelIdx++
	if p.sentinelIdx >= uint8(len(p.sentinelServers)) {
		p.sentinelIdx = 0
	}
	server = p.sentinelServers[p.sentinelIdx]
	info("try connect sentinel server %s", server)
	p.sentinelPool = &RedisPool{}
	p.sentinelPool.Connect(server)
	err = p.sentinelPool.Subscribe("+switch-master")
	if err != nil {
		warn("subscribe sentinel server %s fail, %v", server, err)
		p.sentinelPool = nil
		return
	}
	redisServer = p.sentinelGetMaster()
	if redisServer != "" {
		info("connect to master redis server %s", redisServer)
		p.Connect(redisServer)
	}
	go p.sentinelPool.SubscribeHandle(p.sentinelHandle)
	return
}

func (p *RedisPoolX) sentinelHandle(channel, kind string, data []byte) {
	switch kind {
	case "message":
		if channel == "+switch-master" {
			param := strings.Fields(string(data))
			if len(param) != 5 {
				return
			}
			server := fmt.Sprintf("%s:%s", param[3], param[4])
			info("master redis server switch to %s", server)
			p.Connect(server)
		}
	case "event":
		var err error
		info("sentinel got event %s", string(data))
		err = p.sentinelTryConnect()
		if err != nil {
			warn("try reconnect sentinel fail, %v", err)
		}
	}
}

func (p *RedisPoolX) sentinelGetMaster() (master string) {
	if p.sentinelPool == nil {
		return
	}
	values, err := redis.Strings(p.sentinelPool.Rdo("SENTINEL", "get-master-addr-by-name", p.sentinelMaster))
	if err != nil || len(values) != 2 {
		return
	}
	master = strings.Join(values, ":")
	return
}
