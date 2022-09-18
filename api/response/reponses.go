package response

import (
	"github.com/ditsuke/youtube-focus/internal/yt"
	"github.com/go-chi/render"
	"net/http"
)

type ErrResponse struct {
	Err            error  `json:"-"`
	HttpStatusCode int    `json:"-"`
	StatusText     string `json:"status"`
	ErrorText      string `json:"error,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HttpStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HttpStatusCode: 400,
		StatusText:     "invalid request",
		ErrorText:      err.Error(),
	}
}

type VideosResponse struct {
	Videos []yt.Video `json:"videos"`
	Next   int64      `json:"next"`
}

func NewVideosResponse(videos []yt.Video) *VideosResponse {
	if len(videos) == 0 {
		return &VideosResponse{Videos: videos}
	}

	return &VideosResponse{
		Next:   videos[len(videos)-1].PublishedAt.Unix(),
		Videos: videos,
	}
}

func (v *VideosResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)
	return nil
}
