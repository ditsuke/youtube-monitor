package handlers

import (
	"fmt"
	"github.com/ditsuke/youtube-focus/api/response"
	"github.com/ditsuke/youtube-focus/store"
	"github.com/go-chi/render"
	"net/http"
	"time"
)

const (
	ParamFrom   = "from"
	ParamLimit  = "limit"
	ParamSearch = "search"

	LimitMax     = 20
	LimitDefault = 10

	QueryTimeFmt = time.RFC3339
)

// VideoHandler provides HTTP handlers for the video API.
type VideoHandler struct {
	store store.VideoMetaStore
}

// New returns a VideoHandler configured with the passed store.VideoMetaStore
func New(svc store.VideoMetaStore) *VideoHandler {
	return &VideoHandler{
		store: svc,
	}
}

// Search handles requests with search queries, or without.
func (c *VideoHandler) Search(w http.ResponseWriter, r *http.Request) {
	qParams := r.URL.Query()
	from, limit, err := getPaginationParams(qParams)
	if err != nil {
		_ = render.Render(w, r, response.ErrInvalidRequest(err))
		return
	}

	s, _ := parseParam(qParams, ParamSearch, "")
	if s == "" {
		videos := c.store.Retrieve(from, limit)
		_ = render.Render(w, r, response.NewVideosResponse(videos))
		return
	}

	videos := c.store.Search(s, from, limit)
	_ = render.Render(w, r, response.NewVideosResponse(videos))
}

// AdvancedSearch handles natural-language search queries
func (c *VideoHandler) AdvancedSearch(w http.ResponseWriter, r *http.Request) {
	qParams := r.URL.Query()
	_, limit, err := getPaginationParams(qParams)
	if err != nil {
		_ = render.Render(w, r, response.ErrInvalidRequest(err))
	}

	s, _ := parseParam(qParams, ParamSearch, "")
	if s == "" {
		_ = render.Render(w, r,
			response.ErrInvalidRequest(fmt.Errorf("no `search` parameter in query")))
		return
	}

	videos := c.store.NaturalSearch(s, limit)
	_ = render.Render(w, r, response.NewVideosResponse(videos))
}
