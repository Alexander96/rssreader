# Golang RSS Reader package

## About the package

The goal of the package is to parse and return RSS feed from given array of URLs.

## Usage

```golang
package main

import (
	"fmt"
	rssreader "github.com/Alexander96/rssreader"
)


func main() {
    urls := []string{"http://rss.cnn.com/rss/edition_asia.rss"}
    items, errRss := rssreader.Parse(urls)
    if errRss != nil {
	    fmt.Printf("ERROR: %s\n", errRss)
    }
    fmt.Println(items)
}

```

### Golang version

go1.19.2

