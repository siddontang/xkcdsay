package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
	"github.com/siddontang/xkcdsay"
)

var db *sql.DB

func panicErr(err error) {
	if err == nil {
		return
	}

	panic(err.Error())
}

func getEnvWithDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func init() {
	user := getEnvWithDefault("USER", "root")
	password := getEnvWithDefault("PASS", "")
	host := getEnvWithDefault("HOST", "tidb.5b486b69.1ef404eb.us-west-2.prod.aws.tidbcloud.com")
	port := getEnvWithDefault("PORT", "4000")
	database := getEnvWithDefault("DB", "xkcd")

	dsn := fmt.Sprintf("%s@tcp(%s:%s)/%s", strings.Join([]string{user, password}, ":"),
		host, port, database)
	var err error
	db, err = sql.Open("mysql", dsn)
	panicErr(err)
}

func SyncHandler() (string, error) {
	row := db.QueryRow("select max(xkcd_id) from xkcd")
	panicErr(row.Err())
	var maxID int
	if err := row.Scan(&maxID); err != nil {
		return "", err
	}

	current, err := xkcdsay.GetComicMeta(0)
	if err != nil {
		return "", err
	}

	id := maxID + 1
	if id > current.Num {
		return "", nil
	}

	c, err := xkcdsay.GetComic(id)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://xkcd.com/%d/", c.Num)
	log.Printf("syncing %s", url)

	if err = c.Save(db); err != nil {
		return "", err
	}

	return url, nil
}

func main() {
	lambda.Start(SyncHandler)
}
