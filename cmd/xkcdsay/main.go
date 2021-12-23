package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/base64"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"os"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/BurntSushi/graphics-go/graphics"
	"github.com/mattn/go-sixel"

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

func render(data []byte) error {
	reader := bytes.NewReader(data)
	img, _, err := image.Decode(reader)
	if err != nil {
		return err
	}

	h := img.Bounds().Dy() * 2
	w := img.Bounds().Dx() * 2

	var size [4]uint16
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(os.Stdout.Fd()), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&size)), 0, 0, 0); err != 0 {
		panic(err)
	}
	_, _, width, _ := size[0], size[1], size[2], size[3]

	scale := 1.0
	if t := int(width / 2); t > w {
		scale = float64(t) / float64(w)
	}

	// scale the image to the half of window if possible
	tmp := image.NewNRGBA64(image.Rect(0, 0, int(float64(w)*scale), int(float64(h)*scale)))
	err = graphics.Scale(tmp, img)
	panicErr(err)

	buf := bufio.NewWriter(os.Stdout)
	defer buf.Flush()

	enc := sixel.NewEncoder(buf)
	enc.Dither = true
	return enc.Encode(tmp)
}

func main() {
	flag.Parse()

	rand.Seed(time.Now().Unix())

	dsn := fmt.Sprintf("%s@tcp(%s:%d)/%s?", strings.Join([]string{*user, *password}, ":"),
		*host, *port, *database)
	db, err := sql.Open("mysql", dsn)
	panicErr(err)
	defer db.Close()

	row := db.QueryRow("select count(1) from xkcd")
	panicErr(row.Err())
	var count int
	err = row.Scan(&count)
	panicErr(err)

	id := rand.Intn(count) + 1
	row = db.QueryRow(fmt.Sprintf("select url, file_content from xkcd where xkcd_id = %d;", id))
	panicErr(row.Err())

	var (
		url     string
		content string
	)
	err = row.Scan(&url, &content)
	panicErr(err)

	fmt.Printf("xkcd url: %s\n", url)
	data, err := base64.StdEncoding.DecodeString(content)
	panicErr(err)

	render(data)
}
