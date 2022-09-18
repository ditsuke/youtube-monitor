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
	ParamFrom  = "from"
	ParamLimit = "limit"
	ParamQuery = "q"

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

// Get handles requests for
func (c *VideoHandler) Get(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	marker, limit, err := getPaginationParams(query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
	}

	videos := c.store.Retrieve(marker, limit)

	// j, _ := json.Marshal(videos)
	// _, _ = w.Write(j)
	_ = render.Render(w, r, response.NewVideosResponse(videos))
}

// Search handles requests with search queries
func (c *VideoHandler) Search(w http.ResponseWriter, r *http.Request) {
	qParams := r.URL.Query()
	marker, limit, err := getPaginationParams(qParams)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	s, _ := parseParam(qParams, ParamQuery, "")
	if s == "" {
		_ = render.Render(w, r, response.ErrInvalidRequest(fmt.Errorf("no search query")))
		return
	}

	videos := c.store.Search(s, marker, limit)

	_ = render.Render(w, r, response.NewVideosResponse(videos))
}
