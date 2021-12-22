package xkcdsay

import (
	"bytes"
	"image"
	"math/rand"
	"testing"
)

func checkImage(t *testing.T, c *Comic) {
	_, _, err := image.Decode(bytes.NewReader(c.Content))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDownload(t *testing.T) {
	// get current comic
	c, err := GetComic(0)
	if err != nil {
		t.Fatal(err)
	}
	checkImage(t, c)

	// random get a comic
	num := rand.Intn(c.Num) + 1
	c, err = GetComic(num)
	if err != nil {
		t.Fatal(err)
	}

	if c.Num != num {
		t.Fatalf("want %d but got %d", num, c.Num)
	}
	checkImage(t, c)
}
