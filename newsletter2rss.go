package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

var (
	//AllEmails in memory list of all emails, TODO should put into database
	AllEmails = []Email{}
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

}

func startWebServer() {
	m := martini.Classic()
	m.Use(render.Renderer())

	//API: Create new feed
	m.Post("/feed", func(req *http.Request, r render.Render) {
		log.Printf("post params: %#v", req.Form)
		name := req.PostFormValue("name")
		if !isValidFeedName(name) {
			r.JSON(http.StatusBadRequest, map[string]string{"message": "feed name not found or invalid in request"})
			return
		}

		feed := NewFeed(name)

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
func isValidFeedName(name string) bool {
	return len(name) > 0
}

// NewsLetterFeed is the newsletter which can be subscribed by using an email address
type NewsLetterFeed struct {
	ID         string
	Title      string
	SiteURL    string
	UsedEmails []EmailAddr
	Email      EmailAddr
	Path       string
}

// NewFeed create new feed
func NewFeed(title string) *NewsLetterFeed {
	feed := &NewsLetterFeed{}

	feed.Title = title
	feed.ID = generateFeedID(title)
	feed.Email = NewEmailAddr(feed.ID)

	addEmail(feed.Email, feed)

	return feed
}

//TODO
func generateFeedID(title string) string {
	return title
}

// Email is the data structure represents an email address
type Email struct {
	Addr EmailAddr
	Feed *NewsLetterFeed
}

// EmailAddr an email address
type EmailAddr string

// NewEmailAddr TODO generate a new email address based on feed id
func NewEmailAddr(feedid string) EmailAddr {
	return EmailAddr(fmt.Sprintf("%s@%s", feedid, "localhost"))
}

func addEmail(addr EmailAddr, feed *NewsLetterFeed) {
	AllEmails = append(AllEmails, Email{addr, feed})
}
