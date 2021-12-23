package main

import (
	"encoding/base64"
	"encoding/csv"
	"flag"
	"fmt"
	"os"

	"github.com/siddontang/xkcdsay"
)

var (
	output = flag.String("o", "./xkcd.csv", "output file")
)

func panicErr(err error) {
	if err == nil {
		return
	}

	panic(err.Error())
}

func dumpComic(w *csv.Writer, c *xkcdsay.Comic) {
	item := []string{
		fmt.Sprintf("%d", c.Num),
		c.Title,
		fmt.Sprintf("https://xkcd.com/%d/", c.Num),
		base64.StdEncoding.EncodeToString(c.Content),
		fmt.Sprintf("%s-%s-%s", c.Year, c.Month, c.Day),
		c.Alt,
	}

	err := w.Write(item)
	panicErr(err)
}

func main() {
	current, err := xkcdsay.GetComicMeta(0)
	panicErr(err)

	n := current.Num
	f, err := os.OpenFile(*output, os.O_CREATE|os.O_WRONLY, 0644)
	panicErr(err)
	defer f.Close()

	w := csv.NewWriter(f)
	for i := 1; i <= n; i++ {
		if i == 404 {
			// Aha, I guess you know 404 here
			continue
		}

		c, err := xkcdsay.GetComic(i)
		fmt.Printf("Dump https://xkcd.com/%d/, Img: %s\n", c.Num, c.ImgUrl)
		panicErr(err)
		dumpComic(w, c)
	}

	w.Flush()
	panicErr(w.Error())
}
