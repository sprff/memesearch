package client

import (
	"api-client/pkg/models"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	c := &Client{Url: "http://localhost:1781"}
	ctx := context.Background()
	var id models.MemeID
	var err error
	meme := models.Meme{
		BoardID:      "test",
		Filename:     "some.mp4",
		Descriptions: map[string]string{"sub": "кот"},
	}
	t.Run("Post meme", func(t *testing.T) {
		id, err = c.PostMeme(ctx, meme)
		meme.ID = id
		assert.NoError(t, err)
	})

	t.Run("Get meme", func(t *testing.T) {
		nmeme, err := c.GetMemeByID(ctx, id)
		assert.NoError(t, err)
		nmeme.CreatedAt = time.Time{}
		nmeme.UpdatedAt = time.Time{}
		assert.Equal(t, meme, nmeme)
	})
	t.Run("Put meme", func(t *testing.T) {
		meme.BoardID = "test2"
		err = c.PutMeme(ctx, meme)
		assert.NoError(t, err)
		nmeme, err := c.GetMemeByID(ctx, id)
		assert.NoError(t, err)
		nmeme.CreatedAt = time.Time{}
		nmeme.UpdatedAt = time.Time{}
		assert.Equal(t, meme, nmeme)
	})

}
