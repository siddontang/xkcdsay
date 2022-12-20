package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/siddontang/xkcdsay"
)

var db *sql.DB

func getEnvWithDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func init() {
	user := getEnvWithDefault("USER", xkcdsay.DefaultRoot)
	password := getEnvWithDefault("PASS", "")
	host := getEnvWithDefault("HOST", xkcdsay.DefaultHost)
	port, _ := strconv.Atoi(getEnvWithDefault("PORT", strconv.Itoa(xkcdsay.DefaultPort)))
	database := getEnvWithDefault("DB", xkcdsay.DefaultDB)

	mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: host,
	})

	db = xkcdsay.OpenDB(user, password, host, port, database)
}

func SyncHandler() (string, error) {
	row := db.QueryRow("select max(xkcd_id) from xkcd")
	xkcdsay.PanicErr(row.Err())
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
