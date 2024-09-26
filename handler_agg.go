package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/squashd/blog-aggregator-cli/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.Name)
	}

	interval, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid interval: %w", err)
	}

	log.Printf("Collecting feeds every %s...", interval)

	ticker := time.NewTicker(interval)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}
	log.Printf("Fetching feed %s\n", feed.Name)
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	ctx := context.Background()
	err := db.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s as fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		log.Printf("Couldn't fetch feed %s: %v", feed.Name, err)
		return
	}
	for _, item := range feedData.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err := db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			Title:       item.Title,
			Url:         feed.Url,
			Description: item.Description,
			PublishedAt: publishedAt.Time,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}

	}
	log.Printf("Fetched feed %s, total items %d\n", feed.Name, len(feedData.Channel.Item))
}
