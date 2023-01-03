package rssreader

import (
	"encoding/xml"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_Data(t *testing.T) {

	dataFiles := []string{"cnn_posts1.xml", "cnn_posts2.xml", "cnn_posts3.xml"}

	for _, file := range dataFiles {
		path := fmt.Sprintf("testdata/%s", file)
		f, _ := os.Open(path)
		decoder := xml.NewDecoder(f)

		rss, err := ParseData(decoder)

		if err != nil {
			assert.NotNil(t, err)
			assert.Nil(t, rss)
		} else {
			assert.NotNil(t, rss)
			assert.Nil(t, err)
		}
	}
}
