package db

import (
	"database/sql"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	// blank import for mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// Adapter implements the DbPort interface
type Adapter struct {
	db *sql.DB
}

// NewAdapter creates a new Adapter
func NewAdapter(driverName, dataSourceName string) (*Adapter, error) {
	// connect
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatalf("db connection failur: %v", err)
	}

	// test db connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("db ping failure: %v", err)
	}

	return &Adapter{db: db}, nil
}

// CloseDbConnection closes the db  connection
func (da Adapter) CloseDbConnection() {
	err := da.db.Close()
	if err != nil {
		log.Fatalf("db close failure: %v", err)
	}
}

/*SaveFunction(name string, binary string, cost time.Duration) error
CheckName(name string)error*/

// AddToHistory adds the result of an operation to the database history table
func (da Adapter) SaveFunction(name string, binary string, cost time.Duration) error {
	queryString, args, err := sq.Insert("functions").Columns("name", "location", "cost").
		Values(name, binary, cost).ToSql()
	if err != nil {
		return err
	}

	_, err = da.db.Exec(queryString, args...)
	if err != nil {
		return err
	}

	return nil
}

// CheckName(name string)error*
func (da Adapter) CheckName(name string) error {
	function := sq.Select("*").From("functions")

	active := function.Where(sq.Eq{"name": name})

	queryString, args, err := active.ToSql()
	if err != nil {
		return err
	}

	rows, err := da.db.Query(queryString, args...)
	if err != nil {
		return err
	}
	n := ""
	n2  :=  ""
	n3 := time.Duration(0)
	for rows.Next() {
		rows.Scan(&n, &n2, &n3)
	}
	if n != "" {
		return err
	}
	return nil
}
