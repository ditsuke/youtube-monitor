package interfaces

import "time"

type Store[T any, M time.Time] interface {
	Save([]T)
	Retrieve(marker M, limit int) []T
	Search(query string, marker M, limit int) []T
}
