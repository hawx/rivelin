package models

import (
	"html/template"
	"strings"
	"time"
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
	return template.HTML("<time pubdate=\"" + t.Format(time.RFC3339) + "\">" +
		t.Format("02 Jan; 15:04 PM") +
		"</time>")
}

type River struct {
	UpdatedFeeds UpdatedFeeds
	// Metadata     Metadata
}

func (r River) SetLocation(loc time.Location) River {
	return River{UpdatedFeeds: r.UpdatedFeeds.SetLocation(loc)}
}

type UpdatedFeeds struct {
	UpdatedFeed []UpdatedFeed
}

func (r UpdatedFeeds) SetLocation(loc time.Location) UpdatedFeeds {
	for i, _ := range r.UpdatedFeed {
		r.UpdatedFeed[i] = r.UpdatedFeed[i].SetLocation(loc)
	}
	return r
}

type UpdatedFeed struct {
	FeedUrl         string
	WebsiteUrl      string
	FeedTitle       string
	FeedDescription string
	WhenLastUpdate  RssTime
	Item            []Item
}

func (r UpdatedFeed) SetLocation(loc time.Location) UpdatedFeed {
	r.WhenLastUpdate = RssTime{r.WhenLastUpdate.In(&loc)}
	for i, _ := range r.Item {
		r.Item[i] = r.Item[i].SetLocation(loc)
	}
	return r
}

type Item struct {
	Body      string
	Permalink string
	PubDate   RssTime
	Title     string
	Link      string
	Id        string
}

func (r Item) FilteredBody() string {
	r.Body = strings.TrimSpace(r.Body)

	if strings.HasPrefix(r.Body, "&amp;lt;") ||
		strings.HasPrefix(r.Body, "var gaJsHost") {
		return ""
	}

	return r.Body
}

func (r Item) SetLocation(loc time.Location) Item {
	r.PubDate = RssTime{r.PubDate.In(&loc)}
	return r
}

type Metadata struct {
	Docs      string
	WhenGMT   RssTime
	WhenLocal RssTime
	Version   int
	Secs      int
}
