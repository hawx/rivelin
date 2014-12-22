package models

import (
	"time"
	"html/template"
)

type RssTime struct {
	time.Time
}

func (t RssTime) MarshalText() ([]byte, error) {
	return []byte(t.Format(time.RFC1123Z)), nil
}

func (t RssTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Format(time.RFC1123Z) + `"`), nil
}

func (t *RssTime) UnmarshalText(data []byte) error {
	g, err := time.Parse(time.RFC1123Z, string(data))
	if err != nil {
		return err
	}
	*t = RssTime{g}
	return nil
}

func (t *RssTime) UnmarshalJSON(data []byte) error {
	g, err := time.Parse(`"`+time.RFC1123Z+`"`, string(data))
	if err != nil {
		return err
	}
	*t = RssTime{g}
	return nil
}

func (t *RssTime) HtmlFormat() template.HTML {
	return template.HTML("<time pubdate=\"" + t.Format(time.RFC3339) + "\">" + t.UTC().Format("02 Jan; 15:05 PM") + "</time>")
}

type River struct {
	UpdatedFeeds UpdatedFeeds
	Metadata     Metadata
}

type UpdatedFeeds struct {
	UpdatedFeed []UpdatedFeed
}

type UpdatedFeed struct {
	FeedUrl         string
	WebsiteUrl      string
	FeedTitle       string
	FeedDescription string
	WhenLastUpdate  RssTime
	Item            []Item
}

type Item struct {
	Body      string
	Permalink string
	PubDate   RssTime
	Title     string
	Link      string
	Id        string
}

type Metadata struct {
	Docs      string
	WhenGMT   RssTime
	WhenLocal RssTime
	Version   int
	Secs      int
}
