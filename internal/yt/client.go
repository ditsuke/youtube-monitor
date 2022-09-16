package yt

import (
	"context"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Client struct {
	service *youtube.Service
}

func New(ctx context.Context, apiKey string) (*Client, error) {
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failure: %+v", err)
	}

	return &Client{service}, nil
}

func (c *Client) Query(query string) (*youtube.SearchListResponse, error) {
	return c.service.Search.List([]string{"snippet"}).Type("video").Q(query).Do()
}
