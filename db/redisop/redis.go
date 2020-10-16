package redisop

import (
	"errors"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

var defaultServer string = ":6379"

// common errors
var (
	ErrNoConnection = errors.New("redisop: no connection")
	ErrNoSubscribe  = errors.New("redisop: no subscribe")
)

// PubSubHandle 处理redis订阅消息
type PubSubHandle func(channel, kind string, data []byte)

// RedisPool redis连接池
type RedisPool struct {
	sub *redis.PubSubConn
	rp  *redis.Pool
}

// Connect 连接redis服务器
func (p *RedisPool) Connect(redisServer string) (ok bool) {
	var c redis.Conn
	if redisServer == "default" || redisServer == "" {
		redisServer = defaultServer
	}
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", redisServer)
		return
	}
	p.rp = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 1200 * time.Second,
		Dial:        dialFunc,
	}
	ok = true
	c = p.rp.Get()
	if c.Err() != nil {
		warn("connect to redis server %s fail:%v", redisServer, c.Err())
		ok = false
		p.rp = nil
	}
	defer c.Close()
	return
}

// IsConnected 是否正常连接到服务器
func (p *RedisPool) IsConnected() (ok bool) {
	return p.rp != nil
}

// Close 关闭连接池
func (p *RedisPool) Close() (ok bool) {
	var err error
	ok = true
	if p.rp != nil {
		err = p.rp.Close()
		if err != nil {
			ok = false
		} else {
			p.rp = nil
		}
	}
	return
}

// Rdo 执行redis命令
func (p *RedisPool) Rdo(commandName string, args ...interface{}) (reply interface{}, err error) {
	var c redis.Conn
	if p.rp == nil {
		err = ErrNoConnection
		return
	}
	c = p.rp.Get()
	defer c.Close()
	return c.Do(commandName, args...)
}

// SubscribeHandle 处理redis订阅消息
func (p *RedisPool) SubscribeHandle(handle PubSubHandle) (err error) {
	if p.sub == nil {
		err = ErrNoSubscribe
		return
	}
	for {
		switch v := p.sub.Receive().(type) {
		case redis.Message:
			handle(v.Channel, "message", v.Data)
		case redis.PMessage:
			handle(v.Channel, "pmessage", v.Data)
		case redis.Subscription:
			handle(v.Channel, v.Kind, []byte(fmt.Sprintf("%d", v.Count)))
		case error:
			handle("", "event", []byte(v.Error()))
			return v
		}
	}
}

// UnSubscribe 取消订阅
func (p *RedisPool) UnSubscribe(topic string, isPattern ...bool) (err error) {
	if p.sub == nil {
		return ErrNoSubscribe
	}
	if len(isPattern) > 0 && isPattern[0] == true {
		p.sub.PUnsubscribe(topic)
	} else {
		p.sub.Unsubscribe(topic)
	}
	return
}

// Subscribe 订阅
func (p *RedisPool) Subscribe(topic string, isPattern ...bool) (err error) {
	var c redis.Conn
	if p.rp == nil {
		err = ErrNoConnection
		return
	}
	if p.sub == nil {
		c = p.rp.Get()
		p.sub = &redis.PubSubConn{c}
		if p.sub == nil {
			err = ErrNoSubscribe
			return
		}
	}
	if len(isPattern) > 0 && isPattern[0] == true {
		p.sub.PSubscribe(topic)
	} else {
		p.sub.Subscribe(topic)
	}
	return
}

// GetNotifyConfig 获取redis事件通知配置
func (p *RedisPool) GetNotifyConfig() (cfg string) {
	if vv, err := redis.Strings(p.Rdo("CONFIG", "GET", "notify-keyspace-events")); err != nil || len(vv) < 2 {
		cfg = ""
	} else {
		cfg = vv[1]
	}
	return
}

// SetNotifyConfig 设置redis事件通知配置
func (p *RedisPool) SetNotifyConfig(cfg string) (err error) {
	_, err = p.Rdo("CONFIG", "SET", "notify-keyspace-events", cfg)
	return
}

func defaultPubSubHandle(channel, kind string, data []byte) {
	info("defaultPubSubHandle %s: %s %s\n", channel, kind, string(data))
}
