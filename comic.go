package xkcdsay

import (
	"encoding/json"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
)

// Comic represents a XKCD comic
type Comic struct {
	Num     int    `json:"num"`
	Year    string `json:"year"`
	Month   string `json:"month"`
	Day     string `json:"day"`
	Title   string `json:"title"`
	Alt     string `json:"alt"`
	ImgUrl  string `json:"img"`
	Content []byte `json:"_"`
}

const (
	comicUrl        string = "https://xkcd.com/%d/info.0.json"
	currentComicUrl string = "https://xkcd.com/info.0.json"
)

// GetComicMeta gets a Comic meta from xkcd. If n = 0, gets the current comic
// This function will not down the image, you need to call DownImg later.
func GetComicMeta(n int) (*Comic, error) {
	url := currentComicUrl
	if n > 0 {
		url = fmt.Sprintf(comicUrl, n)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("get meta failed, invalid status code %d", resp.StatusCode)
	}

	var c Comic
	d := json.NewDecoder(resp.Body)
	if err = d.Decode(&c); err != nil {
		return nil, err
	}
	return &c, err
}

// GetComic gets a Comic from xkcd. If n = 0, gets the current comic
func GetComic(n int) (*Comic, error) {
	c, err := GetComicMeta(n)
	if err != nil {
		return nil, err
	}

	if err = c.DownImg(); err != nil {
		return nil, err
	}

	return c, err
}

const noRefImgUrl = "https://imgs.xkcd.com/comics/"

// DownImg downs the image
func (c *Comic) DownImg() error {
	// corner cases
	if c.ImgUrl == noRefImgUrl {
		if c.Num == 1608 {
			// the img link in info json is https://imgs.xkcd.com/comics
			// but the real link is https://xkcd.com/1608/1000:-1074+s.png
			c.ImgUrl = "https://xkcd.com/1608/1000:-1074+s.png"
		}

		// TODO: refine img URL if possible

		// For 1663
		// We can't find the real img URL
		return nil
	}

	resp, err := http.Get(c.ImgUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("get image failed: invalid status code %d", resp.StatusCode)
	}

	c.Content, err = ioutil.ReadAll(resp.Body)
	return err
}
