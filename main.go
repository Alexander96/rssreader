package rssreader

import (
	"fmt"
	"time"
)

type RssItem struct {
	Title       string
	Source      string
	SourceUrl   string
	Link        string
	PublishDate time.Time
	Description string
}

func Parse(urls []string) []RssItem {
	fmt.Println(urls)

	items := make([]RssItem, 0)

	return items
}
