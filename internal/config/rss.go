package config

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type Client struct {
	httpClient http.Client
}

func (s State) FetchFeed(ctx context.Context, fedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("unable to send request: %v", err)
	}

	clnt := Client{
		httpClient: http.Client{},
	}
	req.Header.Add("User-Agent", "gator")

	resp, err := clnt.httpClient.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("response error: %v", err)
	}

	defer resp.Body.Close()

	newRSSFeed := RSSFeed{}
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("read error: %v", err)
	}
	xml.Unmarshal(dat, &newRSSFeed)

	newRSSFeed.Channel.Title = html.UnescapeString(newRSSFeed.Channel.Title)
	newRSSFeed.Channel.Description = html.UnescapeString(newRSSFeed.Channel.Description)
	for i, itm := range newRSSFeed.Channel.Item {
		newRSSFeed.Channel.Item[i].Title = html.UnescapeString(itm.Title)
		newRSSFeed.Channel.Item[i].Description = html.UnescapeString(itm.Description)
	}

	return &newRSSFeed, nil
}
