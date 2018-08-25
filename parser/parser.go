// Package parser provide html email parsers
// this package contains different email parsers for each newsletter mail
// each parser should be ablt to parse a specfic newsletter email's html Contents
// and return a list of articles, which can be conposed into rss feeds
package parser

import (
	"fmt"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

//Article one article
type Article struct {
	Title       string
	Link        string
	Summary     string
	Author      string
	PublishDate time.Time
	Image       string
	Videos      []string
	Domain      string
	Content     string
}

//Parser the common parser interface
type Parser interface {
	Parse(html []byte) ([]Article, error)
}

type FeedItemParser interface {
	Parse(item *gofeed.Item) ([]Article, error)
}

//FindParser find proper parser by the newsletter's information
func FindParser(fromMail string, subject string, html string) (Parser, error) {
	if strings.Contains(fromMail, "@microservicesweekly.com") {
		return &MicroserviceParser{}, nil
	}

	if strings.Contains(fromMail, "@hndigest.com") {
		return &HackerNewsDigestParser{}, nil
	}

	return nil, nil
}

//FindFeedParser find proper html parser for the specified feed
func FindFeedParser(url, title string) (FeedItemParser, error) {
	if strings.Contains(url, "www.discoverdev.io") {
		return &DiscoverdevParser{}, nil
	}
	if strings.Contains(url, "hdchina.org") {
		return &HDCFeedParser{}, nil
	}

	return nil, fmt.Errorf("could not find parser for feed %s", title)
}
