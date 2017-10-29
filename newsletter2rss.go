package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"text/template"
	"time"

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
	//ConfigFilePath path of default config file
	ConfigFilePath = "./config.json"
	//Configuration is the configuration of server
	Configuration struct {
		Name    string `json:"name"`
		DataDir string `json:"datadir"`
		WebDir  string `json:"webdir"`
		Feed    struct {
			Host string `json:"host"`
			Path string `json:"path"`
		} `json:"feed"`
		Email struct {
			ArchiveDir string `json:"archivedir"`
			Host       string `json:"host"`
			Port       string `json:"port"`
			SSL        bool   `json:"ssl"`
			SSLCert    string `json:"sslcert"`
			SSLKey     string `json:"sslkey"`
		} `json:"email"`
	}
)

var feedTmpl = template.Must(template.ParseFiles("./templates/feed.tmpl"))

func main() {
	readConfig()

	loadData()
	startMailServer()
	startWebServer()
}

func readConfig() {
	defaultConfig()

	configFile, err := ioutil.ReadFile(ConfigFilePath)
	if err != nil {
		log.Printf("failed to load config file: %v", err)
		log.Printf("use default configuration")
		return
	}

	err = json.Unmarshal(configFile, &Configuration)
	if err != nil {
		log.Fatalf("failed to parse json config file %s: %v", ConfigFilePath, err)
	}

	log.Printf("Configuration loaded: %+v", Configuration)

	os.MkdirAll(Configuration.DataDir, 0700)
	os.MkdirAll(Configuration.Feed.Path, 0700)
	os.MkdirAll(Configuration.Email.ArchiveDir, 0700)
}

func defaultConfig() {
	Configuration.Name = "news2rss"
	Configuration.DataDir = "./data"
	Configuration.WebDir = "./web-app/build/"
	Configuration.Feed.Host = "localhost"
	Configuration.Feed.Path = "./data/feeds"
	Configuration.Email.ArchiveDir = path.Join(Configuration.DataDir, "mail-archieve")
	Configuration.Email.Host = "localhost"
	Configuration.Email.Port = "2525"
	Configuration.Email.SSL = false
}

func loadData() {
	feeds, err := ioutil.ReadFile("./data/feeds.json")
	if err != nil {
		log.Printf("failed to read feeds.json file, skip silently")
		return
	}

	err = json.Unmarshal(feeds, &AllFeeds)
	if err != nil {
		log.Printf("failed to parse feeds.json, skip silently")
		return
	}

	for i := 0; i < len(AllFeeds); i++ {
		addEmail(AllFeeds[i].Email, &AllFeeds[i])
	}
}

func printEmailIndex() {
	for k := range EmailIndex {
		log.Printf("key: %s, feed: %s", k, EmailIndex[k].Feed.ID)
	}
}

func startMailServer() {
	handler := func(remoteAddr net.Addr, from string, tos []string, data []byte) {
		log.Printf("got mail from %s, remote address: %s", from, remoteAddr)
		log.Printf("recipients: %v", tos)

		archiveMail(from, data)

		emails := findEmails(EmailIndex, tos)
		if len(emails) == 0 {
			log.Printf("could not find any mails match this email, skip")
			return
		}

		log.Printf("found %d feeds match the recipients mail address", len(emails))
		log.Printf("%+v", emails)

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
			log.Printf("update feed: %s", email.Feed.ID)
			email.Feed.Update(entries)
		}
	}

	go func() {
		server := fmt.Sprintf(":%s", Configuration.Email.Port)
		log.Printf("start smtp server on %s", server)

		if Configuration.Email.SSL {
			log.Fatal(smtpd.ListenAndServeTLS(server, Configuration.Email.SSLCert,
				Configuration.Email.SSLKey, handler, Configuration.Name,
				Configuration.Email.Host))
		} else {
			log.Fatal(smtpd.ListenAndServe(server, handler, Configuration.Name,
				Configuration.Email.Host))
		}
	}()
}

func archiveMail(from string, data []byte) {
	log.Printf("save received mail to archive")
	filename := fmt.Sprintf("%s-%d.txt", from, time.Now().Unix())
	filename = path.Join(Configuration.Email.ArchiveDir, filename)

	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		log.Printf("failed to write mail to file: %v", err)
	} else {
		log.Printf("mail achived to file: %s", filename)
	}
}

func findEmails(emailIndex map[string]Email, addr []string) []Email {
	emails := []Email{}
	for _, a := range addr {
		if email, ok := emailIndex[a]; ok {
			log.Printf("found email for feed: %s -> %s", a, email.Feed.ID)
			emails = append(emails, email)
		}
	}
	return emails
}

// NewsLetterFeed is the newsletter which can be subscribed by using an email address
type NewsLetterFeed struct {
	ID         string      `json:"id"`
	Title      string      `json:"title"`
	URL        string      `json:"url"`
	Updated    string      `json:"-"`
	UsedEmails []string    `json:"usedEmails"`
	Email      string      `json:"email"`
	Path       string      `json:"-"`
	Entries    []FeedEntry `json:"-"`
}

//FeedEntry feed entry
type FeedEntry struct {
	ID      string
	Updated string
	Title   string
	Summary string
}

//Update update entries for a feed
func (feed *NewsLetterFeed) Update(entries []FeedEntry) error {
	feed.Entries = entries

	feed.Updated = time.Now().Format(time.RFC3339)

	feedFilePath := "./data/feeds/" + feed.ID + ".xml"

	feedFile, err := os.OpenFile(feedFilePath, os.O_WRONLY|os.O_CREATE, 0755)
	defer feedFile.Close()
	if err != nil {
		log.Printf("failed to open feed file to write: %v", err)
		return err
	}

	err = feedTmpl.Execute(feedFile, feed)
	if err != nil {
		log.Printf("failed to execute feed template: %v", err)
		return err
	}

	log.Printf("updated feed with %d entries", len(entries))
	return nil
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

	feed.URL = fmt.Sprintf("/feeds/%s.xml", feed.ID)
	feed.Email = mail

	return feed
}

func convertArticleToEntry(articles []parser.Article) []FeedEntry {
	entries := []FeedEntry{}

	for _, article := range articles {
		entry := FeedEntry{
			ID:      article.Link,
			Title:   article.Title,
			Summary: article.Summary,
			Updated: time.Now().Format(time.RFC3339),
		}
		entries = append(entries, entry)
	}

	return entries
}

func startWebServer() {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(martini.Static(Configuration.WebDir,
		martini.StaticOptions{Exclude: "/feeds/"}))

	//API: Create new feed
	m.Post("/api/feeds", func(req *http.Request, r render.Render) {
		decoder := json.NewDecoder(req.Body)
		data := struct {
			Title string `json:"title"`
		}{}

		err := decoder.Decode(&data)

		if err != nil {
			log.Printf("failed to parse request body as json: %v", err)
			r.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request data"})
			return
		}

		if !isValidFeedTitle(data.Title) {
			r.JSON(http.StatusBadRequest, map[string]string{"message": "feed title not found or invalid in request"})
			return
		}

		mail := newMailAddr(data.Title)
		feed := NewFeed(data.Title, mail)
		addEmail(feed.Email, feed)
		AllFeeds = append(AllFeeds, *feed)

		saveFeedsData()

		r.JSON(200, feed)
	})

	m.Get("/api/feeds", func(r render.Render) {
		r.JSON(http.StatusOK, AllFeeds)
	})

	m.Get("/feeds/:feed", func(params martini.Params, r render.Render) {
		feedFilename := fmt.Sprintf("%s/%s", Configuration.Feed.Path, params["feed"])

		feedData, err := ioutil.ReadFile(feedFilename)
		if err != nil {
			log.Printf("failed to read feed file %s: %v", feedFilename, err)
			if os.IsNotExist(err) {
				r.Error(http.StatusNotFound)
			} else {
				r.Error(http.StatusInternalServerError)
			}
			return
		}

		log.Printf("response with feed file data: %s", feedFilename)
		r.Header().Set(render.ContentType, "application/atom+xml")
		r.Data(http.StatusOK, feedData)
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
		addr := fmt.Sprintf("%s@%s", result, Configuration.Email.Host)
		_, ok := EmailIndex[addr]
		return ok
	}

	return fmt.Sprintf("%s@%s", slug.New(title, exist), Configuration.Email.Host)
}

func addEmail(addr string, feed *NewsLetterFeed) {
	if exist, ok := EmailIndex[addr]; ok {
		log.Printf("adding existing email %s, exist feed: %s, new feed: %s. Skip add",
			addr, exist.Feed.Title, feed.Title)
		return
	}
	EmailIndex[addr] = Email{addr, feed}
}

func saveFeedsData() {
	data, err := json.Marshal(AllFeeds)
	if err != nil {
		log.Printf("failed to marshal feeds to json: %v", err)
		return
	}

	err = ioutil.WriteFile("./data/feeds.json", data, 0644)

	if err != nil {
		log.Printf("failed to write feeds to json file: %v", err)
	}
}
