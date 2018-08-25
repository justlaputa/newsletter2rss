package parser

import (
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

//DiscoverdevParser mail parser for microservice weekly newsletter
type DiscoverdevParser struct {
}

func (dp *DiscoverdevParser) Parse(item *gofeed.Item) ([]Article, error) {
	articles := []Article{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(item.Description))
	if err != nil {
		log.Printf("failed to create goquery doc, %v", err)
		return articles, err
	}

	articlesElem := doc.Find("ul>li>a")

	articlesElem.Each(func(i int, elem *goquery.Selection) {
		href, exist := elem.Attr("href")
		if !exist {
			log.Printf("no link found in this element, %d, %v", i, elem)
			return
		}
		title := elem.Text()

		domain := ""
		linkURL, err := url.Parse(href)
		if err != nil {
			log.Printf("url is not valid: %s", href)
		} else {
			domain = linkURL.Hostname()
		}

		articles = append(articles, Article{Title: title, Link: href, Summary: "", Domain: domain})
	})

	return articles, nil
}
