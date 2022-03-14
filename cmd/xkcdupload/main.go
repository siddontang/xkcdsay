package main

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var (
	host      = flag.String("H", "tidb.b68f76dc.1ef404eb.us-west-2.prod.aws.tidbcloud.com", "Host")
	port      = flag.Int("P", 4000, "Port")
	user      = flag.String("u", "guest", "user")
	password  = flag.String("pass", "11111111", "password")
	database  = flag.String("D", "xkcd", "database")
	inputFile = flag.String("I", "", "CSV file will be uploaded")
	batchSize = flag.Int("S", 100, "Number of comics to be inserted in one batch")
)

func panicErr(err error) {
	if err == nil {
		return
	}

	panic(err.Error())
}

func uploadComics(db *sql.DB, records [][]string) {
	var buf bytes.Buffer

	args := make([]interface{}, 0, 6*len(records))
	for i, record := range records {
		buf.WriteString(`(`)
		for j, r := range record {
			buf.WriteString(`?`)
			args = append(args, r)
			if j != len(record)-1 {
				buf.WriteString(`,`)
			}
		}

		if i == len(records)-1 {
			buf.WriteString(`)`)
		} else {
			buf.WriteString(`),`)
		}
	}

	query := fmt.Sprintf("replace into xkcd (xkcd_id, title, url, file_content, date_published, alt) values %s", buf.String())
	_, err := db.Exec(query, args...)
	panicErr(err)

	fmt.Printf("insert %d comics\n", len(records))
}

func main() {
	flag.Parse()

	dsn := fmt.Sprintf("%s@tcp(%s:%d)/%s?maxAllowedPacket=0", strings.Join([]string{*user, *password}, ":"),
		*host, *port, *database)
	db, err := sql.Open("mysql", dsn)
	panicErr(err)
	defer db.Close()

	f, err := os.Open(*inputFile)
	panicErr(err)
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ','
	r.LazyQuotes = true

	records := make([][]string, 0, *batchSize)

	for {
		records = records[0:0]

		for i := 0; i < *batchSize; i++ {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			panicErr(err)

			records = append(records, record)
		}

		uploadComics(db, records)

		if len(records) < *batchSize {
			break
		}
	}
}
