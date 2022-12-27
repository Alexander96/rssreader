package rssreader

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

type customTime struct {
	time.Time
}

func (c *customTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parse, err := http.ParseTime(v)
	if err != nil {
		return err
	}
	*c = customTime{parse}
	return nil
}

type rss struct {
	Channel channel `xml:"channel"`
}

type channel struct {
	Title         string     `xml:"title"`
	Link          string     `xml:"link"`
	Description   string     `xml:"description"`
	Items         []RssItem  `xml:"item"`
	LastBuildDate customTime `xml:"lastBuildDate"`
}

type RssItem struct {
	Title       string     `xml:"title"`
	Link        string     `xml:"link"`
	GUID        string     `xml:"guid"`
	PublishDate customTime `xml:"pubDate"`
	Description string     `xml:"description"`
}

func Parse(urls []string) ([]RssItem, error) {
	var items []RssItem
	for _, url := range urls {
		output, err := parseUrl(url)
		if err != nil {
			fmt.Printf("ERROR while parsing url %s\n", url)
			return items, err
		}
		items = append(items, output...)
	}

	return items, nil
}

func parseUrl(url string) ([]RssItem, error) {
	var items []RssItem

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("ERROR GET url %s\n", url)
		return items, err
	}
	defer resp.Body.Close()

	rss := rss{}

	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&rss)
	if err != nil {
		fmt.Printf("ERROR Decode: %v\n", err)
		return items, err
	}

	for _, item := range rss.Channel.Items {
		items = append(items, item)
	}

	return items, nil
}
