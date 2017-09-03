package main

import "log"

// NewsLetterFeed is the newsletter which can be subscribed by using an email address
type NewsLetterFeed struct {
	ID         string
	Title      string
	SiteURL    string
	UsedEmails []string
	Email      string
	Path       string
}

type FeedEntry struct {
}

func (feed *NewsLetterFeed) Update(entries []FeedEntry) {
	log.Printf("updating %d entries", len(entries))
}

// NewFeed create new feed
func NewFeed(title string) *NewsLetterFeed {
	feed := &NewsLetterFeed{}

	feed.Title = title
	feed.ID = generateFeedID(title)
	feed.Email = NewEmailAddr(feed.ID)

	return feed
}

//TODO
func generateFeedID(title string) string {
	return title
}
