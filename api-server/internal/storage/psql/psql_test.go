package psql

import (
	"context"
	"memesearch/internal/config"
	"memesearch/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getConfig() config.DatabaseConfig {
	//TODO use env_vars
	return config.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "user",
		Password: "password",
		Dbname:   "meme-search-test",
	}
}

func TestBoard(t *testing.T) {
	cfg := getConfig()
	ctx := context.Background()

	store, err := NewBoardStore(cfg)
	require.NoError(t, err)
	board := models.Board{
		Owner: "test_owner",
		Name:  "test",
	}
	var id models.BoardID

	t.Run("Insert board", func(t *testing.T) {
		b, err := store.CreateBoard(ctx, "test_owner", "test")
		id = b.ID
		require.NoError(t, err)
		board.ID = id
	})

	t.Run("Get board", func(t *testing.T) {
		nboard, err := store.GetBoardByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, board, nboard)
	})

	t.Run("Get board not exist", func(t *testing.T) {
		_, err := store.GetBoardByID(ctx, "unknown_id")
		require.Equal(t, models.ErrBoardNotFound, err)

	})

	t.Run("Update board", func(t *testing.T) {
		board.Name = "TEEST"
		err := store.UpdateBoard(ctx, board)
		require.NoError(t, err)
		nboard, err := store.GetBoardByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, board, nboard)
	})
	t.Run("Update board not exist", func(t *testing.T) {
		err := store.UpdateBoard(ctx, models.Board{ID: "unknonwn_id"})
		assert.Equal(t, models.ErrBoardNotFound, err)
	})

	t.Run("Delete board", func(t *testing.T) {
		err := store.DeleteBoard(ctx, board.ID)
		assert.NoError(t, err)
		err = store.DeleteBoard(ctx, "unknonwn_id")
		assert.Equal(t, models.ErrBoardNotFound, err)
	})
}

func TestMedia(t *testing.T) {
	cfg := getConfig()
	ctx := context.Background()

	store, err := NewMediaStore(cfg)
	require.NoError(t, err)

	t.Run("Get", func(t *testing.T) {
		_, err := store.GetMediaByID(ctx, "_test_id")
		require.NoError(t, err)
	})

	t.Run("Set", func(t *testing.T) {
		media, err := store.GetMediaByID(ctx, "_test_id")
		require.NoError(t, err)
		media.ID = "_test_id2"
		err = store.SetMediaByID(ctx, media)
		require.NoError(t, err)
		nmedia, err := store.GetMediaByID(ctx, "_test_id2")
		require.NoError(t, err)
		assert.Equal(t, media.Body, nmedia.Body)
	})

	t.Run("Clean", func(t *testing.T) {
		err = store.SetMediaByID(ctx, models.Media{ID: "_test_id2", Body: []byte{}})
		require.NoError(t, err)
	})
}

func TestMeme(t *testing.T) {
	cfg := getConfig()
	ctx := context.Background()

	store, err := NewMemeStore(cfg)
	require.NoError(t, err)
	meme := models.Meme{
		BoardID:  "_test_board",
		Filename: "file.mp4",
		Description: map[string]string{
			"subject": "кот",
			"text":    "я кот",
		},
	}

	var id models.MemeID
	t.Run("Insert meme", func(t *testing.T) {
		id, err = store.InsertMeme(ctx, meme)
		require.NoError(t, err)
		meme.ID = id

	})
	t.Run("Get meme", func(t *testing.T) {
		nmeme, err := store.GetMemeByID(ctx, id)
		require.NoError(t, err)
		nmeme.CreatedAt = time.Time{}
		nmeme.UpdatedAt = time.Time{}
		assert.Equal(t, meme, nmeme)
	})
	t.Run("Get meme not exist", func(t *testing.T) {
		_, err := store.GetMemeByID(ctx, "unknown_id")
		require.Equal(t, models.ErrMemeNotFound, err)

	})
	t.Run("Update meme", func(t *testing.T) {
		meme.Filename = "TEEST.png"
		err := store.UpdateMeme(ctx, meme)
		require.NoError(t, err)
		nmeme, err := store.GetMemeByID(ctx, id)
		require.NoError(t, err)
		nmeme.CreatedAt = time.Time{}
		nmeme.UpdatedAt = time.Time{}
		assert.Equal(t, meme, nmeme)
	})
	t.Run("Update meme not exist", func(t *testing.T) {
		err := store.UpdateMeme(ctx, models.Meme{ID: "unknonwn_id"})
		assert.Equal(t, models.ErrMemeNotFound, err)
	})
	t.Run("Delete meme", func(t *testing.T) {
		err := store.DeleteMeme(ctx, meme.ID)
		assert.NoError(t, err)
		err = store.DeleteMeme(ctx, "unknonwn_id")
		assert.Equal(t, models.ErrMemeNotFound, err)
	})
}
func TestUser(t *testing.T) {
	// TODO
}
