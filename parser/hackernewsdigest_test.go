package parser

import (
	"log"
	"os"
	"testing"

	"github.com/jhillyerd/enmime"
)

func TestHackerNewsDigestParser(t *testing.T) {
	r, err := os.Open("sample_email_tests/hackernewsdigest.txt")
	if err != nil {
		t.Error(err)
	}

	env, err := enmime.ReadEnvelope(r)
	if err != nil {
		t.Error(err)
	}
	log.Printf("HTML Body: %v chars\n", len(env.HTML))

	hn := &HackerNewsDigestParser{}

	articles, err := hn.Parse([]byte(env.HTML))
	if err != nil {
		t.Errorf("failed to parse mail from hacker news digest, %v", err)
	}

	for _, a := range articles {
		log.Printf("Title %s", a.Title)
		log.Printf("Link: <%s>", a.Link)
		log.Printf("Summary: %s", a.Summary)
	}
}
