package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/jhillyerd/enmime"
	"github.com/justlaputa/newsletter2rss/parser"
	"github.com/justlaputa/newsletter2rss/slug"
	"github.com/martini-contrib/render"
	"github.com/mhale/smtpd"
)

var (
	//EmailIndex in memory index of all emails, TODO should put into database
	EmailIndex = make(map[string]Email)
	//AllFeeds in memory list of all feeds, TODO should put into database
	AllFeeds = []NewsLetterFeed{}
)

func main() {
	readConfig()

	startMailServer()
	startWebServer()
}

func readConfig() {

}

func startMailServer() {
	handler := func(remoteAddr net.Addr, from string, tos []string, data []byte) {
		log.Printf("got mail from %s, remote address: %s", from, remoteAddr)
		log.Printf("recipients: %v", tos)

		emails := findEmails(EmailIndex, tos)
		if len(emails) == 0 {
			log.Printf("could not find any mails match this email, skip")
			return
		}

		log.Printf("found %d feeds match the recipients mail address", len(emails))

		message, err := enmime.ReadEnvelope(bytes.NewReader(data))
		if err != nil {
			log.Printf("failed to parse email, skip: %v", err)
			return
		}

		mailParser := parser.FindParser(message.GetHeader("From"), message.GetHeader("Subject"), message.HTML)
		articles, err := mailParser.Parse([]byte(message.HTML))
		if err != nil {
			log.Printf("failed to parse email content, %v", err)
			return
		}

		if len(articles) == 0 {
			log.Printf("no articles found in the email, skip")
		}

		log.Printf("found %d articles in the email", len(articles))

		entries := convertArticleToEntry(articles)

		for _, email := range emails {
			email.Feed.Update(entries)
		}
	}

	go func() {
		log.Printf("start smtp server on :2525")
		log.Fatal(smtpd.ListenAndServe(":2525", handler, "news2rss", "localhost"))
	}()
}

func findEmails(emailIndex map[string]Email, addr []string) []Email {
	emails := []Email{}
	for _, a := range addr {
		if email, ok := emailIndex[a]; ok {
			emails = append(emails, email)
		}
	}
	return emails
}

// NewsLetterFeed is the newsletter which can be subscribed by using an email address
type NewsLetterFeed struct {
	ID         string
	Title      string
	SiteURL    string
	UsedEmails []string
	Email      string
	Path       string
}

//FeedEntry feed entry
type FeedEntry struct {
}

//Update update entries for a feed
func (feed *NewsLetterFeed) Update(entries []FeedEntry) {
	log.Printf("updating %d entries", len(entries))
}

// NewFeed create new feed
func NewFeed(title, mail string) *NewsLetterFeed {
	feed := &NewsLetterFeed{}

	feed.Title = title
	feed.ID = slug.New(title, func(id string) bool {
		for _, f := range AllFeeds {
			if f.ID == id {
				return true
			}
		}
		return false
	})

	feed.Email = mail

	return feed
}

func convertArticleToEntry(articles []parser.Article) []FeedEntry {
	return []FeedEntry{}
}

func startWebServer() {
	m := martini.Classic()
	m.Use(render.Renderer())

	//API: Create new feed
	m.Post("/feed", func(req *http.Request, r render.Render) {
		log.Printf("post params: %#v", req.Form)
		title := req.PostFormValue("title")
		if !isValidFeedTitle(title) {
			r.JSON(http.StatusBadRequest, map[string]string{"message": "feed title not found or invalid in request"})
			return
		}

		mail := newMailAddr(title)
		feed := NewFeed(title, mail)
		addEmail(feed.Email, feed)
		AllFeeds = append(AllFeeds, *feed)

		r.JSON(200, map[string]string{"id": feed.ID, "email": string(feed.Email)})
	})

	//Page: home page, list all feeds
	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", "hello world")
	})

	//TODO run on another port
	m.Run()
}

//TODO
func isValidFeedTitle(name string) bool {
	return len(name) > 0
}

// Email is the data structure represents an email address
type Email struct {
	Addr string
	Feed *NewsLetterFeed
}

func newMailAddr(title string) string {
	exist := func(result string) bool {
		addr := result + "@localhost"
		_, ok := EmailIndex[addr]
		return ok
	}

	return fmt.Sprintf("%s@%s", slug.New(title, exist), "localhost")
}

func addEmail(addr string, feed *NewsLetterFeed) {
	if exist, ok := EmailIndex[addr]; ok {
		log.Printf("adding existing email %s, exist feed: %s, new feed: %s. Skip add", addr, exist.Feed.Title, feed.Title)
		return
	}

	EmailIndex[addr] = Email{addr, feed}
}
