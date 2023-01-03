package rssreader

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"sync"
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

func (r RssItem) String() string {
	return fmt.Sprintf("Title: %s\nLink: %s\nGUID: %s\nPublishData: %s\nDescription: %s\n", r.Title, r.Link, r.GUID, r.PublishDate, r.Description)
}

func Parse(urls []string) ([]RssItem, error) {
	var items []RssItem
	var chans []chan RssItem
	var wg sync.WaitGroup
	wg.Add(len(urls))
	for i := 0; i < len(urls); i++ {
		chans = append(chans, make(chan RssItem))
	}

	for i, url := range urls {
		go parseUrl(url, chans[i], &wg)
	}
	agg := make(chan RssItem)
	for _, ch := range chans {
		go func(c chan RssItem) {
			for msg := range c {
				agg <- msg
			}
			wg.Done()
		}(ch)
	}
	go func() {
		wg.Wait()
		close(agg)
	}()

	for item := range agg {
		items = append(items, item)
	}

	return items, nil
}

func parseUrl(url string, ch chan<- RssItem, wg *sync.WaitGroup) {
	defer close(ch)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("ERROR GET url %s\n", url)
	}
	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)

	rss, err := ParseData(decoder)
	if err != nil {
		fmt.Printf("ERROR Decode: %v\n", err)
	}

	for _, item := range rss.Channel.Items {
		ch <- item
	}
}

func ParseData(decoder *xml.Decoder) (rss, error) {
	rss := rss{}
	err := decoder.Decode(&rss)
	return rss, err
}
