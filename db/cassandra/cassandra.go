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

// Session 返回连接Session
func (a *Agent) Session() *gocql.Session {
	return a.session
}

// Batch 批量操作
func (a *Agent) Batch(typ gocql.BatchType) *gocql.Batch {
	return a.session.NewBatch(typ)
}

// Tables 查询所有表格及其字段名
func (a *Agent) Tables() (tbls map[string][]string) {
	ks, err := a.session.KeyspaceMetadata(a.DefaultKeyspace)
	if err != nil {
		return
	}
	tbls = make(map[string][]string)
	for _, tbl := range ks.Tables {
		tbls[tbl.Name] = tbl.OrderedColumns
	}
	return
}
