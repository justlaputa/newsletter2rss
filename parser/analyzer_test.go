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

func TestMakeContent(t *testing.T) {
	content := "hello\r\n\r\n\"world\"\ntest\r\n\r\n"
	image := "https://test.com/test.jpg"
	link := "https://test.com/article.html"

	expected := "From: <a href=\"https://test.com/article.html\">test.com</a><br/><img src=\"https://test.com/test.jpg\" /><br/>hello<br/><br/>\"world\"<br/>test<br/><br/>"

	result := makeContent(link, image, content)
	resultUnescaped := html.UnescapeString(result)

	if resultUnescaped != expected {
		log.Printf("input:\n%s", content)
		log.Printf("result  : %s", result)
		log.Printf("  result unescape: %s", html.UnescapeString(result))
		log.Printf("expected unescape: %s", expected)
		t.Fatalf("escape not success")
	}
}
