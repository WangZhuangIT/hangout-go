package output

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/kshvakov/clickhouse"
)

type DB struct {
	dbConn *sql.DB
}

func NewDBConn(host string, port string, database string) *DB {
	var err error
	db := new(DB)
	dsn := fmt.Sprintf("tcp://%s:%s?database=%s", host, port, database)
	db.dbConn, err = sql.Open("clickhouse", dsn)
	if err != nil {
		panic(err)
	}
	err = db.dbConn.Ping()
	if err != nil {
		panic(err)
	}
	db.dbConn.SetMaxOpenConns(10)

	return db
}

func (cDb *DB) Insert(table string, param map[string]interface{}) error {
	l := len(param)
	if l <= 0 {
		return errors.New("empty data")
	}

	var columns, holder string
	var values []interface{}
	i := 1
	sql := "INSERT INTO " + table
	for col, val := range param {
		columns += col
		holder += "?"
		values = append(values, val)
		if i < l {
			columns += ", "
			holder += ", "
		}
		i++
	}

	sql += " (" + columns + ") VALUES (" + holder + ")"
	// fmt.Println(sql, values)
	tx, err := cDb.dbConn.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
