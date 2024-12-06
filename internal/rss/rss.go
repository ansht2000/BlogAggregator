package rss

import (
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
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

func stripHTMLTags(input string) string {
	tokenizer := html.NewTokenizer(strings.NewReader(input))
	var buffer bytes.Buffer

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			break
		}
		if tokenType == html.TextToken {
			buffer.WriteString(string(tokenizer.Text()))
		}
	}
	return buffer.String()
}

func unescapeHTMLString(feed *RSSFeed) {
	feed.Channel.Title = stripHTMLTags(html.UnescapeString(feed.Channel.Title))
	feed.Channel.Description = stripHTMLTags(html.UnescapeString(feed.Channel.Description))
	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = stripHTMLTags(html.UnescapeString(item.Title))
		feed.Channel.Item[i].Description = stripHTMLTags(html.UnescapeString(item.Description))
	}
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	req.Header.Set("User-Agent", "gator")
	
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	var feed RSSFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return &RSSFeed{}, err
	}

	unescapeHTMLString(&feed)

	return &feed, nil
}
