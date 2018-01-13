package parser

import (
	"fmt"
	"html"
	"log"
	"net/url"
	"strings"
	"sync"

	textapi "github.com/AYLIEN/aylien_textapi_go"
)

const (
	// AylienAPIMaxConcurrentCount maximum concurrent api request count for aylien
	AylienAPIMaxConcurrentCount = 5
)

// Analyzer interface to provider analyze function for articles
type Analyzer interface {
	Analyze([]Article)
}

// AylienAnalyzer analyzer using aylien api
type AylienAnalyzer struct {
	AppID      string
	AppKey     string
	Client     *textapi.Client
	APIChannel chan struct{}
}

func (ay *AylienAnalyzer) init() error {
	auth := textapi.Auth{ay.AppID, ay.AppKey}
	client, err := textapi.NewClient(auth, true)
	if err != nil {
		log.Printf("could not create aylien api client, %v", err)
		return err
	}
	ay.Client = client
	ay.APIChannel = make(chan struct{}, AylienAPIMaxConcurrentCount)
	for i := 0; i < AylienAPIMaxConcurrentCount; i++ {
		ay.APIChannel <- struct{}{}
	}
	return nil
}

// Analyze use aylien api to extract each article and get summary of them
func (ay *AylienAnalyzer) Analyze(articles []Article) {
	if len(articles) == 0 {
		return
	}

	totalCount := len(articles)

	wg := sync.WaitGroup{}
	wg.Add(totalCount)

	for i := 0; i < totalCount; i++ {
		<-ay.APIChannel

		go func(index int) {
			analyze(ay.Client, &articles[index])

			ay.APIChannel <- struct{}{}
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func analyze(client *textapi.Client, article *Article) {
	log.Printf("use aylien extract api to analyze article")

	extractParams := &textapi.ExtractParams{
		URL:       article.Link,
		BestImage: true,
	}

	result, err := client.Extract(extractParams)
	if err != nil {
		log.Printf("got error while calling aylien api: %v", err)
		return
	}

	article.Author = result.Author
	article.Image = result.Image
	article.PublishDate = result.PublishDate.Time
	article.Videos = result.Videos[:]
	article.Content = makeContent(article.Link, article.Image, result.Article)
}

func makeContent(link, image, contents string) string {
	host := ""
	linkURL, err := url.Parse(link)
	if err != nil {
		log.Printf("failed to parse title url %s: %v", link, err)
	} else {
		host = linkURL.Hostname()
	}
	from := fmt.Sprintf("From: <a href=\"%s\">%s</a><br/>", link, host)
	imageLink := fmt.Sprintf("<img src=\"%s\" /><br/>", image)
	contents = strings.Replace(contents, "\r\n", "\n", -1)
	contents = strings.Replace(contents, "\n", "<br/>", -1)
	contents = from + imageLink + contents

	return html.EscapeString(contents)
}

func summarize(client *textapi.Client, article *Article) {
	log.Printf("getting summary for article: %s", article.Title)
	summarizeParams := &textapi.SummarizeParams{
		URL:               article.Link,
		NumberOfSentences: 3,
		Title:             article.Title,
	}
	summary, err := client.Summarize(summarizeParams)
	if err != nil {
		log.Printf("failed to get summary, %v", err)
		return
	}
	if len(summary.Sentences) != 0 {
		article.Summary = strings.Join(summary.Sentences, "\n")
	}
}

func wrapContent(content string) string {
	paragraphs := strings.Split(content, "\r\n\r\n")
	for i, p := range paragraphs {
		paragraphs[i] = "<p>" + p + "</p>"
	}

	result := strings.Join(paragraphs, "<br>")
	return html.EscapeString(result)
}

// NewAylienAnalyzer create an aylien analyzer
func NewAylienAnalyzer(appid, appkey string) Analyzer {
	analyzer := &AylienAnalyzer{AppID: appid, AppKey: appkey}
	err := analyzer.init()

	if err != nil {
		log.Printf("failed to create aylien api client")
		return nil
	}

	return analyzer
}
