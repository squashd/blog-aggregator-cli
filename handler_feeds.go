package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/squashd/blog-aggregator-cli/internal/database"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "text/xml")
	req.Header.Add("user-agent", "blog-aggregator-cli")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var feed *RSSFeed
	err = xml.NewDecoder(res.Body).Decode(&feed)
	if err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return feed, nil

}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("agg requires a feed URL")
	}
	url := cmd.Args[0]

	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("unable to fetch feed: %w", err)
	}

	fmt.Printf("%+w\n", *feed)

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: %s <feed_name> <feed_url>", cmd.Name)
	}
	ctx := context.Background()

	feedName := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      feedName,
		Url:       url,
		UserID:    user.ID,
	})

	_, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feed.ID,
		UserID:    feed.UserID,
	})
	if err != nil {
		return fmt.Errorf("unable to follow feed: %w", err)
	}

	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("%s\n", feed.Url)
		fmt.Printf("%s\n", feed.Name)
		fmt.Printf("%s\n", feed.UserName)
	}

	return nil
}
