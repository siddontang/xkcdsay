package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"strings"

	"github.com/siddontang/xkcdsay"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	host     = flag.String("H", xkcdsay.DefaultHost, "Host")
	port     = flag.Int("P", xkcdsay.DefaultPort, "Port")
	user     = flag.String("u", xkcdsay.DefaultGuest, "user")
	password = flag.String("pass", xkcdsay.DefaultPass, "password")
	database = flag.String("D", xkcdsay.DefaultDB, "database")
)

func panicErr(err error) {
	if err == nil {
		return
	}

	panic(err.Error())
}

func main() {
	flag.Parse()

	mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: *host,
	})

	dsn := fmt.Sprintf("%s@tcp(%s:%d)/%s?tls=tidb", strings.Join([]string{*user, *password}, ":"),
		*host, *port, *database)
	db, err := sql.Open("mysql", dsn)
	panicErr(err)
	defer db.Close()

	row := db.QueryRow("select max(xkcd_id) from xkcd")
	var maxID sql.NullInt64
	err = row.Scan(&maxID)
	panicErr(err)

	current, err := xkcdsay.GetComicMeta(0)
	panicErr(err)

	mid := 0
	if maxID.Valid {
		mid = int(maxID.Int64)
	}
	for id := mid + 1; id <= current.Num; id++ {
		c, err := xkcdsay.GetComic(id)
		panicErr(err)

		fmt.Printf("syncing https://xkcd.com/%d/\n", c.Num)

		err = c.Save(db)
		panicErr(err)
	}
}
