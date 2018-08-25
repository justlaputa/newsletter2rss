package parser

import (
	"os"
	"testing"

	"github.com/mmcdole/gofeed"
)

func TestHDCFeedParser(t *testing.T) {
	feedFile, err := os.Open("sample_feeds/hdc-feed.xml")
	if err != nil {
		t.Error(err)
	}

	fp := gofeed.NewParser()
	feeds, err := fp.Parse(feedFile)
	if err != nil {
		t.Error(err)
	}

	parser := &HDCFeedParser{}

	articles, err := parser.Parse(feeds.Items[0])

	if len(articles) != 1 {
		t.Fatalf("expect 1 article from 1st feed item, got %d", len(articles))
	}

	expected := Article{
		Title: "2036 Origin Unknown",
		Link:  "https://hdchina.org/details.php?id=290329",
	}
	actual := articles[0]

	if expected.Title != actual.Title ||
		expected.Link != actual.Link {
		t.Fatalf("Expected: %v\nActual: %v", expected, articles[0])
	}
}
