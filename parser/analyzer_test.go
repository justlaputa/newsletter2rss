package parser

import (
	"fmt"
	"testing"
)

func TestAylienAnalyzer(t *testing.T) {
	analyzer := NewAylienAnalyzer("api", "key")

	articles := []Article{
		{
			Link:  "https://open.nytimes.com/https-open-nytimes-com-the-new-york-times-as-a-tor-onion-service-e0d0b67b7482",
			Title: "The New York Times is Now Available as a Tor Onion Service",
		},
	}

	analyzer.Analyze(articles)

	fmt.Printf("%#v", articles)
}
