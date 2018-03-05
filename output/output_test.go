package output

import (
	// "github.com/stretchr/testify/assert"
	"testing"
)

var dbConn *DB

func TestMain(m *testing.M) {
	dbConn = NewDBConn("10.97.14.111", "9000", "test")
	m.Run()
}

func Test_Insert(t *testing.T) {
	param := make(map[string]interface{})
	param["date"] = "2018-02-07"
	param["number"] = 123
	param["data"] = "I'm lihua"
	err := dbConn.Insert("test_plus", param)
	if err != nil {
		t.Error(err)
	}
}
