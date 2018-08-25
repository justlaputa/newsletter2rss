package parser

import (
	"fmt"

	"github.com/justlaputa/movie/pt"
	"github.com/mmcdole/gofeed"
)

// HDCFeedParser parse hdc feed
type HDCFeedParser struct{}

func (p *HDCFeedParser) Parse(item *gofeed.Item) ([]Article, error) {
	article := Article{}

	info := pt.ParseHDCTitle(item.Title)

	article.Title = fmt.Sprintf("%s (%d)", info.Title, info.Year)
	article.Link = item.Link
	article.Summary = fmt.Sprintf("%s %s %dGB", info.Source, info.Resolution, info.Size/1000000000)
	return []Article{article}, nil
}
