package parser

import (
	"fmt"
	"html"
	"log"
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

func TestWrapContent(t *testing.T) {
	content := "a\r\n\r\nb\r\n\r\nc\"\r\n\r\n"
	expected := "&lt;p&gt;a&lt;/p&gt;&lt;br&gt;&lt;p&gt;b&lt;/p&gt;&lt;br&gt;&lt;p&gt;c&#34;&lt;/p&gt;&lt;br&gt;&lt;p&gt;&lt;/p&gt;"

	result := wrapContent(content)

	if result != expected {
		log.Printf("input: %s", content)
		log.Printf("expected: %s", expected)
		log.Printf("result  : %s", result)
		log.Printf("result unescape: %s", html.UnescapeString(result))
		t.Fatalf("escape not success")
	}
}
