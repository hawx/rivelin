package models

import (
	"strings"
	"time"
)

type River struct {
	UpdatedFeeds UpdatedFeeds
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
