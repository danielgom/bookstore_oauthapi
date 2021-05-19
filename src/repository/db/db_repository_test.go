package db

import (
	"github.com/gocql/gocql"
	"testing"
)

// MAYBE A TODO MOCKING QUERIES
type MockQuery func(stmt string, values ...interface{}) *gocql.Query

type MockSession struct {
	MockQuery MockQuery
}

func (m *MockSession) Query(stmt string, values ...interface{}) *gocql.Query {
	return m.MockQuery(stmt, values)
}

func TestRepositoryGetByID(t *testing.T) {

}
