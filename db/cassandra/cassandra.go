package cassandra

import (
	"github.com/gocql/gocql"
)

// Agent 客户端
type Agent struct {
	DefaultKeyspace string
	Cluster         *gocql.ClusterConfig
	session         *gocql.Session
}

// Connect 初始连接
func (a *Agent) Connect(keySpace string, servers ...string) (err error) {
	a.DefaultKeyspace = keySpace
	a.Cluster = gocql.NewCluster(servers...)
	a.Cluster.Keyspace = keySpace
	a.session, err = a.Cluster.CreateSession()
	return
}

// Close 关闭连接
func (a *Agent) Close() {
	if a.session == nil || a.session.Closed() {
		return
	}
	a.session.Close()
}

// Query 查询
func (a *Agent) Query(sql string, values ...interface{}) *gocql.Query {
	if a.session == nil || a.session.Closed() {
		return nil
	}
	return a.session.Query(sql, values...)
}
