package yt

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"net/http"
	"time"
)

type Client struct {
	service *youtube.Service
	logger  zerolog.Logger
}

type Opt func(c *Client)

func WithLogger(logger zerolog.Logger) Opt {
	return func(c *Client) {
		c.logger = logger
	}
}

// New returns a new instance of Client, returns a non-nil error if an error is returned by youtube.NewService
// This function employs the options-pattern to configure the client in-situ.
func New(ctx context.Context, apiKey string, opts ...Opt) (*Client, error) {
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failure: %+v", err)
	}
	client := &Client{service: service}
	for _, opt := range opts {
		opt(client)
	}
	return client, nil
}

// QueryLatestVideos returns Video metas matching a query in reverse chronological order (ie: latest).
// Query is capped at the default 5 items at the moment.
func (c *Client) QueryLatestVideos(query string) ([]Video, error) {
	r, err := c.service.Search.List([]string{"snippet"}).Type("video").Q(query).Order("date").PublishedAfter(time.Now().Add(-1 * 24 * time.Hour).Format(time.RFC3339)).Do()
	if err != nil {
		apiError, ok := err.(*googleapi.Error)
		if !ok {
			return nil, fmt.Errorf("list query: %+v", err)
		}
		if apiError.Code == http.StatusForbidden {
			// Invalid token or rate limit exceeded
			// @todo token switching logic goes here ;)
		}
		return nil, fmt.Errorf("youtube: %+v", err)
	}

	videos := make([]Video, len(r.Items))
	for i, result := range r.Items {
		publish, _ := time.Parse(time.RFC3339, result.Snippet.PublishedAt)
		videos[i] = Video{
			Title:        result.Snippet.Title,
			Description:  result.Snippet.Description,
			VideoId:      result.Id.VideoId,
			PublishedAt:  publish,
			ThumbnailUrl: result.Snippet.Thumbnails.Default.Url,
		}
	}

	return videos, nil
}
