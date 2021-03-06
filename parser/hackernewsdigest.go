package parser

import (
	"bytes"
	"log"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

//MicroserviceParser mail parser for microservice weekly newsletter
type HackerNewsDigestParser struct {
}

//Parse parst html contents and return an article list
func (hn *HackerNewsDigestParser) Parse(html []byte) ([]Article, error) {
	articles := []Article{}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		log.Printf("failed to create goquery doc, %v", err)
		return articles, err
	}

	articlesElem := doc.Find("table>tbody>tr>th>a")

	articlesElem.Each(func(i int, elem *goquery.Selection) {
		href, exist := elem.Attr("href")
		if !exist {
			log.Printf("no link found in this article, %d, %v", i, elem)
			return
		}
		title := elem.Text()
		summary := elem.Siblings().Text()
		domain := ""
		linkURL, err := url.Parse(href)
		if err != nil {
			log.Printf("url is not valid: %s", href)
		} else {
			domain = linkURL.Hostname()
		}

		articles = append(articles, Article{Title: title, Link: href, Summary: summary, Domain: domain})
	})

	return articles, nil
}
