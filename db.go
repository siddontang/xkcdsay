package xkcdsay

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

const (
	DefaultHost  = "gateway01.us-west-2.prod.aws.tidbcloud.com"
	DefaultPort  = 4000
	DefaultRoot  = "4A7D3bbkQWsWSEH.root"
	DefaultGuest = "4A7D3bbkQWsWSEH.guest"
	DefaultPass  = "11111111"
	DefaultDB    = "xkcd"
)

// OpenDB opens the database connection.
func OpenDB(user, password, host string, port int, database string) *sql.DB {
	mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: host,
	})

	dsn := fmt.Sprintf("%s@tcp(%s:%d)/%s?tls=tidb", strings.Join([]string{user, password}, ":"),
		host, port, database)
	db, err := sql.Open("mysql", dsn)
	PanicErr(err)
	return db
}

// PanicErr panics if the error is not nil.
func PanicErr(err error) {
	if err == nil {
		return
	}

	panic(err.Error())
}
