// Package parser provide html email parsers
// this package contains different email parsers for each newsletter mail
// each parser should be ablt to parse a specfic newsletter email's html Contents
// and return a list of articles, which can be conposed into rss feeds
package parser

import "strings"

//Article one article
type Article struct {
	Title   string
	Link    string
	Summary string
}

//Parser the common parser interface
type Parser interface {
	Parse(html []byte) ([]Article, error)
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
