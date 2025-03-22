package httpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"memesearch/internal/config"
	"memesearch/internal/models"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getConfig() config.Config {
	return config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "user",
			Password: "password",
			Dbname:   "meme-search-test",
		},
	}
}

func getTestServer(t *testing.T) *httptest.Server {
	cfg := getConfig()
	server, err := New(cfg)
	require.NoError(t, err)
	return httptest.NewServer(server)
}

func TestServerIsRunning(t *testing.T) {
	ts := getTestServer(t)

	res, err := http.Get(ts.URL + "/about")
	assert.NoError(t, err)
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.Equal(t, "t.me/sprff_code", string(body))
}

func TestServer(t *testing.T) {
	ts := getTestServer(t)
	var id models.MemeID
	meme := models.Meme{
		BoardID:      "board1",
		Filename:     "file.mp4",
		Descriptions: map[string]string{"subject": "кот", "text": "я кот"},
	}
	t.Run("POST /memes", func(t *testing.T) {
		_, resp, err := makeJSONReqest("POST", fmt.Sprintf("%s/memes", ts.URL), meme)
		assert.NoError(t, err)
		assert.Equal(t, "OK", resp["status"])
		data := resp["data"].(map[string]any)
		id = models.MemeID(data["id"].(string))
	})
	t.Run("GET /memes/{id}", func(t *testing.T) {
		_, resp, err := makeJSONReqest("GET", fmt.Sprintf("%s/memes/%s", ts.URL, id), nil)
		assert.NoError(t, err)
		assert.Equal(t, "OK", resp["status"])
		data := resp["data"].(map[string]any)
		assert.Equal(t, data["board_id"].(string), string(meme.BoardID))
		assert.Equal(t, data["filename"].(string), meme.Filename)
		assert.Equal(t, data["descriptions"].(map[string]any), map[string]any{"subject": "кот", "text": "я кот"})
	})
	t.Run("PUT /memes/{id}", func(t *testing.T) {
		meme.BoardID = "board2"
		t.Run("Set", func(t *testing.T) {
			_, resp, err := makeJSONReqest("PUT", fmt.Sprintf("%s/memes/%s", ts.URL, id), meme)
			assert.NoError(t, err)
			assert.Equal(t, "OK", resp["status"])
			data := resp["data"].(map[string]any)
			id = models.MemeID(data["id"].(string))

		})
		t.Run("Get", func(t *testing.T) {
			_, resp, err := makeJSONReqest("GET", fmt.Sprintf("%s/memes/%s", ts.URL, id), nil)
			assert.NoError(t, err)
			assert.Equal(t, "OK", resp["status"])
			data := resp["data"].(map[string]any)
			assert.Equal(t, data["board_id"].(string), string(meme.BoardID))
			assert.Equal(t, data["filename"].(string), meme.Filename)
			assert.Equal(t, data["descriptions"].(map[string]any), map[string]any{"subject": "кот", "text": "я кот"})
		})
	})

	t.Run("PUT /media/{id}", func(t *testing.T) {
		if id == "" {
			id = "test"
		}
		url := fmt.Sprintf("%s/media/%s", ts.URL, id)
		t.Logf("url: %s", url)

		var requestBody bytes.Buffer
		contentType := ""
		{ // prepare body for FormFIle read
			writer := multipart.NewWriter(&requestBody)
			part, err := writer.CreateFormFile("media", "filename.ext")
			if err != nil {
				fmt.Println("Ошибка создания части файла:", err)
				return
			}
			part.Write([]byte("Hello"))
			writer.Close()
			contentType = writer.FormDataContentType()
		}

		request, err := http.NewRequest("PUT", url, &requestBody)
		require.NoError(t, err)
		request.Header.Set("Content-Type", contentType)

		resp, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		m := map[string]any{}
		err = json.Unmarshal(respBody, &m)
		require.NoError(t, err, string(respBody))
		assert.Equal(t, "OK", m["status"].(string))
	})

	t.Run("GET /media/{id}", func(t *testing.T) {
		if id == "" {
			id = "test"
		}
		url := fmt.Sprintf("%s/media/%s", ts.URL, id)
		t.Logf("url: %s", url)

		request, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(request)
		require.NoError(t, err)
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, "Hello", string(respBody))
	})
}

func TestServer_invalidInput(t *testing.T) {
	ts := getTestServer(t)
	t.Run("POST /memes", func(t *testing.T) {
		t.Run("invalid body", func(t *testing.T) {
			_, resp, err := makeJSONReqest("POST", fmt.Sprintf("%s/memes", ts.URL), nil)
			assert.NoError(t, err)
			require.Equal(t, "INVALID_INPUT", resp["status"])
			data := resp["err_data"].(map[string]any)
			assert.Equal(t, "body", data["reason"].(string))
		})
		t.Run("wrong type of field", func(t *testing.T) {
			_, resp, err := makeJSONReqest("POST", fmt.Sprintf("%s/memes", ts.URL), map[string]any{"filename": 15})
			assert.NoError(t, err)
			require.Equal(t, "INVALID_INPUT", resp["status"])
			data := resp["err_data"].(map[string]any)
			assert.Equal(t, "filename expected to be string", data["reason"].(string))
		})
	})

	t.Run("PUT /memes", func(t *testing.T) {
		t.Run("not found", func(t *testing.T) {
			_, resp, err := makeJSONReqest("PUT", fmt.Sprintf("%s/memes/%s", ts.URL, "unknown_id"), map[string]any{})
			assert.NoError(t, err)
			require.Equal(t, "MEME_NOT_FOUND", resp["status"])
		})
		var id models.MemeID
		t.Run("paste empty meme", func(t *testing.T) {
			_, resp, err := makeJSONReqest("POST", fmt.Sprintf("%s/memes", ts.URL), map[string]any{})
			assert.NoError(t, err)
			require.Equal(t, "OK", resp["status"])
			data := resp["data"].(map[string]any)
			id = models.MemeID(data["id"].(string))
		})

		t.Run("invalid body", func(t *testing.T) {
			_, resp, err := makeJSONReqest("PUT", fmt.Sprintf("%s/memes/%s", ts.URL, id), nil)
			assert.NoError(t, err)
			require.Equal(t, "INVALID_INPUT", resp["status"])
			data := resp["err_data"].(map[string]any)
			assert.Equal(t, "body", data["reason"].(string))
		})
		t.Run("wrong type of field", func(t *testing.T) {
			_, resp, err := makeJSONReqest("PUT", fmt.Sprintf("%s/memes/%s", ts.URL, id), map[string]any{"filename": 15})
			assert.NoError(t, err)
			require.Equal(t, "INVALID_INPUT", resp["status"])
			data := resp["err_data"].(map[string]any)
			assert.Equal(t, "filename expected to be string", data["reason"].(string))
		})
	})
	t.Run("GET /memes", func(t *testing.T) {
		t.Run("not found", func(t *testing.T) {
			_, resp, err := makeJSONReqest("GET", fmt.Sprintf("%s/memes/%s", ts.URL, "unknown_id"), map[string]any{})
			assert.NoError(t, err)
			require.Equal(t, "MEME_NOT_FOUND", resp["status"])
		})
	})

	t.Run("GET /media", func(t *testing.T) {
		t.Run("not found", func(t *testing.T) {
			_, resp, err := makeJSONReqest("GET", fmt.Sprintf("%s/media/%s", ts.URL, "unknown_id"), map[string]any{})
			assert.NoError(t, err)
			require.Equal(t, "MEDIA_NOT_FOUND", resp["status"])
		})
	})

	t.Run("PUT /media", func(t *testing.T) {
		t.Run("media file is required", func(t *testing.T) {
			_, resp, err := makeJSONReqest("PUT", fmt.Sprintf("%s/media/%s", ts.URL, "some_id"), map[string]any{})
			assert.NoError(t, err)
			require.Equal(t, "MEDIA_IS_REQUIRED", resp["status"])
		})
	})
}

func makeJSONReqest(method string, url string, body any) (int, map[string]any, error) {
	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return 0, nil, fmt.Errorf("can't marshal body: %w", err)
		}
	}
	bodyReader := bytes.NewBuffer(bodyBytes)
	request, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return 0, nil, fmt.Errorf("can't create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, nil, fmt.Errorf("can't do request: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("can't read body: %w", err)
	}

	m := map[string]any{}
	err = json.Unmarshal(bodyBytes, &m)
	if err != nil {
		return 0, nil, fmt.Errorf("can't unmarshal body: %w", err)
	}

	return resp.StatusCode, m, nil
}
