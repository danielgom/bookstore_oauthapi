package cassandra

import (
	"github.com/danielgom/bookstore_oauthapi/src/repository/db"
	"github.com/gocql/gocql"
)

func Init() {

	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "oauth"
	cluster.Consistency = gocql.Quorum

	var err error
	if db.Session, err = cluster.CreateSession(); err != nil {
		panic(err)
	}
}
