package telegram

import "context"

type CachedMediaType string

const (
	CMVideo CachedMediaType = "video"
	CMPhoto CachedMediaType = "photo"
)

type CachedMedia struct {
	FileID string
	Type   CachedMediaType
}

type CallbackEntry struct {
	Text     string
	Callback func()
}

type MediaGroupEntry struct {
	Media   CachedMedia
	Caption string
}

type CachedMediaStorage interface {
	Get(ctx context.Context, key string) (CachedMedia, error)
	Set(ctx context.Context, key string, value CachedMedia) error
}

type UploadEntry struct {
	Name string //filename without extension
	Body *[]byte
}
