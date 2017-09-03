package parser

import (
	"bytes"
	"log"

	"github.com/PuerkitoBio/goquery"
)

//MicroserviceParser mail parser for microservice weekly newsletter
type MicroserviceParser struct {
}

//Parse parst html contents and return an article list
func (mp *MicroserviceParser) Parse(html []byte) ([]Article, error) {
	articles := []Article{}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		log.Printf("failed to create goquery doc, %v", err)
		return articles, err
	}

	articlesElem := doc.Find("#articles").Parent().Find("h3")

	articlesElem.Each(func(i int, elem *goquery.Selection) {
		href, exist := elem.Find("a").Attr("href")
		if !exist {
			log.Printf("no link found in this article, %d, %v", i, elem)
			return
		}
		title := elem.Find("a").Text()
		summary := elem.NextFiltered("p").Text()

		articles = append(articles, Article{title, href, summary})
	})

	return articles, nil
}
