package main

import (
	"database/sql"
	"flag"
	"fmt"
	"strings"

	"github.com/siddontang/xkcdsay"

	_ "github.com/go-sql-driver/mysql"
)

var (
	host     = flag.String("H", "tidb.5b486b69.1ef404eb.us-west-2.prod.aws.tidbcloud.com", "Host")
	port     = flag.Int("P", 4000, "Port")
	user     = flag.String("u", "guest", "user")
	password = flag.String("pass", "11111111", "password")
	database = flag.String("D", "xkcd", "database")
)

func panicErr(err error) {
	if err == nil {
		return
	}

	panic(err.Error())
}

func main() {
	flag.Parse()

	dsn := fmt.Sprintf("%s@tcp(%s:%d)/%s", strings.Join([]string{*user, *password}, ":"),
		*host, *port, *database)
	db, err := sql.Open("mysql", dsn)
	panicErr(err)
	defer db.Close()

	row := db.QueryRow("select max(xkcd_id) from xkcd")
	panicErr(row.Err())
	var maxID int
	err = row.Scan(&maxID)
	panicErr(err)

	current, err := xkcdsay.GetComicMeta(0)
	panicErr(err)

	for id := maxID + 1; id <= current.Num; id++ {
		c, err := xkcdsay.GetComic(id)
		panicErr(err)

		fmt.Printf("syncing https://xkcd.com/%d/", c.Num)

		err = c.Save(db)
		panicErr(err)
	}
}
