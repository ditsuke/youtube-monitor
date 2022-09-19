package yt

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"net/http"
	"sync"
	"time"
)

// Client is a wrapper around the YouTube Data API v3, with a method to fetch the latest videos
// along with multi-token support to circumvent rate-limiting.
type Client struct {
	logger       zerolog.Logger
	tokens       []string
	muTokenState struct {
		sync.Mutex
		Ptr        int
		lastWorked bool
	}
}

type Opt func(c *Client)

func WithLogger(logger zerolog.Logger) Opt {
	return func(c *Client) {
		c.logger = logger
	}
}

// New returns a new instance of Client, returns a non-nil error if an error is returned by youtube.NewService
// This function employs the options pattern to configure the client in-situ.
func New(apiKeys []string, opts ...Opt) (*Client, error) {
	if len(apiKeys) == 0 {
		return nil, errors.New("no api keys!")
	}
	client := &Client{tokens: apiKeys}
	client.muTokenState.lastWorked = true
	for _, opt := range opts {
		opt(client)
	}
	return client, nil
}

// QueryLatestVideos returns Video metas matching a query in reverse chronological order (ie: latest).
// Query is capped at the default 5 items at the moment.
func (c *Client) QueryLatestVideos(query string) ([]Video, error) {
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(c.getCurrentToken()))
	if err != nil {
		return nil, errors.Wrap(err, "youtube api service")
	}

	r, err := service.Search.List([]string{"snippet"}).Type("video").Q(query).Order("date").PublishedAfter(time.Now().Add(-1 * 24 * time.Hour).Format(time.RFC3339)).Do()
	if err != nil {
		apiError, ok := err.(*googleapi.Error)
		if !ok {
			return nil, fmt.Errorf("list query: %+v", err)
		}
		if apiError.Code == http.StatusForbidden {
			// Invalid token or rate limit exceeded
			c.logger.Warn().Int("apiKeyIndex", c.muTokenState.Ptr).Msg("current api key exhausted")
			// @todo account for invalid api keys
			if ok := c.useNextToken(false); ok {
				return c.QueryLatestVideos(query)
			}
			err := fmt.Errorf("youtube tokens exhausted")
			c.logger.Error().Err(err).Msg("exiting youtube poller")
			return nil, err
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

	// employ the round-robin strategy to cycle between tokens
	c.useNextToken(true)
	return videos, nil
}

// getCurrentToken returns the API token currently in-use by the client
func (c *Client) getCurrentToken() string {
	c.muTokenState.Lock()
	defer c.muTokenState.Unlock()
	return c.tokens[c.muTokenState.Ptr]
}

// useNextToken switches out the API token in-use by the client to bypass token-based rate limiting.
func (c *Client) useNextToken(currentWorked bool) bool {
	c.muTokenState.Lock()
	defer c.muTokenState.Unlock()
	if !c.muTokenState.lastWorked && !currentWorked {
		c.logger.Error().Msg("last api token failed too. consider adding more tokens")
		return false
	}
	c.muTokenState.Ptr = (c.muTokenState.Ptr + 1) % len(c.tokens)
	c.muTokenState.lastWorked = false
	c.logger.Debug().Int("index", c.muTokenState.Ptr).Msg("switching youtube api token")
	return true
}
