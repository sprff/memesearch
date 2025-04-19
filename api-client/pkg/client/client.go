package client

import (
	"api-client/internal/apiclient"
	"api-client/pkg/models"
	"bytes"
	"context"
	"fmt"
	"net/http"
)

var _ ClientInterface = Client{}

type Client struct {
	api       *apiclient.ClientWithResponses
	token     string
	requestId string
}

type ClientInterface interface {
	About(ctx context.Context) (info map[string]any, err error)
	AuthRegister(ctx context.Context, login, password string) (id models.UserID, err error)
	AuthLogin(ctx context.Context, login, password string) (token string, err error)
	AuthWhoami(ctx context.Context) (user models.User, err error)
	ListBoards(ctx context.Context, page, pageSize int, sortBy string) (boards []models.Board, err error)
	PostBoard(ctx context.Context, name string) (board models.Board, err error)
	DeleteBoardByID(ctx context.Context, boardID models.BoardID) (board models.Board, err error)
	GetBoardByID(ctx context.Context, boardID models.BoardID) (board models.Board, err error)
	UpdateBoardByID(ctx context.Context, boardID models.BoardID, name *string, owner *models.UserID) (board models.Board, err error)
	GetMediaByID(ctx context.Context, mediaID models.MediaID) (media models.Media, err error)
	PutMediaByID(ctx context.Context, media models.Media) (err error)
	ListMemes(ctx context.Context, page, pageSize int, sortBy string) (boards []models.Meme, err error)
	PostMeme(ctx context.Context, boardID models.BoardID, filename string, dsc map[string]string) (meme models.Meme, err error)
	DeleteMemeByID(ctx context.Context, memeID models.MemeID) (meme models.Meme, err error)
	GetMemeByID(ctx context.Context, memeID models.MemeID) (meme models.Meme, err error)
	UpdateMemeByID(ctx context.Context, memeID models.MemeID, boardID *models.BoardID, filename *string, dsc *map[string]string) (meme models.Meme, err error)
	SearchMemes(ctx context.Context, page, pageSize int, general string) (memes []models.ScoredMeme, err error)
	SubscribeByBoardID(ctx context.Context, boardID models.BoardID) (err error)
	UnsubscribeByBoardID(ctx context.Context, boardID models.BoardID) (err error)
	GetUserByID(ctx context.Context, userID models.UserID) (user models.User, err error)
}

func New(url string) (Client, error) {
	c, err := apiclient.NewClientWithResponses(url)
	if err != nil {
		return Client{}, fmt.Errorf("can't create client with responses: %w", err)
	}
	return Client{api: c}, nil
}

// About implements ClientInterface.
func (c Client) About(ctx context.Context) (info map[string]any, err error) {
	info = make(map[string]any)
	resp, err := c.api.AboutWithResponse(ctx, c.middlewares()...)
	if err != nil {
		err = fmt.Errorf("can't request: %w", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		info["apiname"] = resp.JSON200.ApiName
		info["version"] = resp.JSON200.Version
		return
	case 401:
		err = models.ErrUnauthorized
		return
	case 403:
		err = models.ErrForbidden
		return
	default:
		return nil, fmt.Errorf("unexpected response %d: %s", resp.StatusCode(), string(resp.Body))
	}
}

// AuthLogin implements ClientInterface.
func (c Client) AuthLogin(ctx context.Context, login string, password string) (token string, err error) {
	resp, err := c.api.AuthLoginWithResponse(ctx, apiclient.AuthLoginJSONRequestBody{Login: login, Password: password}, c.middlewares()...)
	if err != nil {
		err = fmt.Errorf("can't request: %w", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		token = resp.JSON200.Token
		return
	case 401:
		err = models.ErrUnauthorized
		return
	case 403:
		err = models.ErrForbidden
		return
	case 404:
		err = models.ErrUserNotFound
		return
	case 400:
		err = parseApiError(*resp.JSON400)
		return
	default:
		err = fmt.Errorf("unexpected response %d: %s", resp.StatusCode(), string(resp.Body))
		return
	}
}

// AuthRegister implements ClientInterface.
func (c Client) AuthRegister(ctx context.Context, login string, password string) (id models.UserID, err error) {
	resp, err := c.api.AuthRegisterWithResponse(ctx, apiclient.AuthRegisterJSONRequestBody{Login: login, Password: password}, c.middlewares()...)
	if err != nil {
		err = fmt.Errorf("can't request: %w", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		id = models.UserID(resp.JSON200.Id)
		return
	case 400:
		err = parseApiError(*resp.JSON400)
		return
	case 401:
		err = models.ErrUnauthorized
		return
	case 403:
		err = models.ErrForbidden
		return
	case 409:
		err = models.ErrLoginExists
		return
	default:
		err = fmt.Errorf("unexpected response %d: %s", resp.StatusCode(), string(resp.Body))
		return
	}
}

// AuthWhoami implements ClientInterface.
func (c Client) AuthWhoami(ctx context.Context) (user models.User, err error) {
	resp, err := c.api.AuthWhoamiWithResponse(ctx, c.middlewares()...)
	if err != nil {
		err = fmt.Errorf("can't request: %w", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		user = convertUserToModel(*resp.JSON200)
		return
	case 401:
		err = models.ErrUnauthorized
		return
	case 403:
		err = models.ErrForbidden
		return
	default:
		err = fmt.Errorf("unexpected response %d: %s", resp.StatusCode(), string(resp.Body))
		return
	}
}

// DeleteBoardByID implements ClientInterface.
func (c Client) DeleteBoardByID(ctx context.Context, boardID models.BoardID) (board models.Board, err error) {
	resp, err := c.api.DeleteBoardByIDWithResponse(ctx, apiclient.BoardId(boardID), c.middlewares()...)
	if err != nil {
		err = fmt.Errorf("can't request: %w", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		board = convertBoardToModel(*resp.JSON200)
		return
	case 401:
		err = models.ErrUnauthorized
		return
	case 403:
		err = models.ErrForbidden
		return
	case 404:
		err = models.ErrBoardNotFound
		return
	default:
		err = fmt.Errorf("unexpected response %d: %s", resp.StatusCode(), string(resp.Body))
		return
	}
}

// DeleteMemeByID implements ClientInterface.
func (c Client) DeleteMemeByID(ctx context.Context, memeID models.MemeID) (meme models.Meme, err error) {
	resp, err := c.api.DeleteMemeByIDWithResponse(ctx, apiclient.MemeId(memeID), c.middlewares()...)
	if err != nil {
		err = fmt.Errorf("can't request: %w", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		// TODO
		// meme = convertMemeToModel(*resp.JSON200)
		return
	case 401:
		err = models.ErrUnauthorized
		return
	case 403:
		err = models.ErrForbidden
		return
	case 404:
		err = models.ErrMemeNotFound
		return
	default:
		err = fmt.Errorf("unexpected response %d: %s", resp.StatusCode(), string(resp.Body))
		return
	}
}

// GetBoardByID implements ClientInterface.
func (c Client) GetBoardByID(ctx context.Context, boardID models.BoardID) (board models.Board, err error) {
	resp, err := c.api.GetBoardByIDWithResponse(ctx, apiclient.BoardId(boardID), c.middlewares()...)
	if err != nil {
		err = fmt.Errorf("can't request: %w", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		board = convertBoardToModel(*resp.JSON200)
		return
	case 401:
		err = models.ErrUnauthorized
		return
	case 403:
		err = models.ErrForbidden
		return
	case 404:
		err = models.ErrBoardNotFound
		return
	default:
		err = fmt.Errorf("unexpected response %d: %s", resp.StatusCode(), string(resp.Body))
		return
	}
}

// GetMediaByID implements ClientInterface.
func (c Client) GetMediaByID(ctx context.Context, mediaID models.MediaID) (media models.Media, err error) {
	resp, err := c.api.GetMediaByIDWithResponse(ctx, apiclient.MediaId(mediaID), c.middlewares()...)
	if err != nil {
		err = fmt.Errorf("can't request: %w", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		media.Body = resp.Body
		media.ID = mediaID
		return
	case 401:
		err = models.ErrUnauthorized
		return
	case 403:
		err = models.ErrForbidden
		return
	case 404:
		err = models.ErrMediaNotFound
		return
	default:
		err = fmt.Errorf("unexpected response %d: %s", resp.StatusCode(), string(resp.Body))
		return
	}
}

// GetMemeByID implements ClientInterface.
func (c Client) GetMemeByID(ctx context.Context, memeID models.MemeID) (meme models.Meme, err error) {
	resp, err := c.api.GetMemeByIDWithResponse(ctx, apiclient.MemeId(memeID), c.middlewares()...)
	if err != nil {
		err = fmt.Errorf("can't request: %w", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		meme = convertMemeToModel(*resp.JSON200)
		return
	case 401:
		err = models.ErrUnauthorized
		return
	case 403:
		err = models.ErrForbidden
		return
	case 404:
		err = models.ErrMemeNotFound
		return
	default:
		err = fmt.Errorf("unexpected response %d: %s", resp.StatusCode(), string(resp.Body))
		return
	}
}

// GetUserByID implements ClientInterface.
func (c Client) GetUserByID(ctx context.Context, userID models.UserID) (user models.User, err error) {
	resp, err := c.api.GetUserByIDWithResponse(ctx, apiclient.UserId(userID), c.middlewares()...)
	if err != nil {
		err = fmt.Errorf("can't request: %w", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		user = convertUserToModel(*resp.JSON200)
		return
	case 401:
		err = models.ErrUnauthorized
		return
	case 403:
		err = models.ErrForbidden
		return
	case 404:
		err = models.ErrUserNotFound
		return
	default:
		err = fmt.Errorf("unexpected response %d: %s", resp.StatusCode(), string(resp.Body))
		return
	}
}

// ListBoards implements ClientInterface.
func (c Client) ListBoards(ctx context.Context, page int, pageSize int, sortBy string) (boards []models.Board, err error) {
	req := &apiclient.ListBoardsParams{Page: &page, PageSize: &pageSize, SortBy: (*apiclient.ListBoardsParamsSortBy)(&sortBy)}
	resp, err := c.api.ListBoardsWithResponse(ctx, req, c.middlewares()...)
	if err != nil {
		err = fmt.Errorf("can't request: %w", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		//TODO unify have(ListMemes) and don't have(ListBoards) items
		for _, b := range *resp.JSON200 {
			boards = append(boards, convertBoardToModel(b))
		}
		return
	case 401:
		err = models.ErrUnauthorized
		return
	case 403:
		err = models.ErrForbidden
		return
	default:
		err = fmt.Errorf("unexpected response %d: %s", resp.StatusCode(), string(resp.Body))
		return
	}
}

// ListMemes implements ClientInterface.
func (c Client) ListMemes(ctx context.Context, page int, pageSize int, sortBy string) (memes []models.Meme, err error) {
	req := &apiclient.ListMemesParams{Page: &page, PageSize: &pageSize, SortBy: (*apiclient.ListMemesParamsSortBy)(&sortBy)}
	resp, err := c.api.ListMemesWithResponse(ctx, req, c.middlewares()...)
	if err != nil {
		err = fmt.Errorf("can't request: %w", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		for _, m := range resp.JSON200.Items {
			memes = append(memes, convertMemeToModel(m))
		}
		return
	case 401:
		err = models.ErrUnauthorized
		return
	case 403:
		err = models.ErrForbidden
		return
	default:
		err = fmt.Errorf("unexpected response %d: %s", resp.StatusCode(), string(resp.Body))
		return
	}
}

// PostBoard implements ClientInterface.
func (c Client) PostBoard(ctx context.Context, name string) (board models.Board, err error) {
	req := apiclient.PostBoardJSONRequestBody{Name: name}
	resp, err := c.api.PostBoardWithResponse(ctx, req, c.middlewares()...)
	if err != nil {
		err = fmt.Errorf("can't request: %w", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		board = convertBoardToModel(*resp.JSON200)
		return
	case 401:
		err = models.ErrUnauthorized
		return
	case 403:
		err = models.ErrForbidden
		return
	default:
		err = fmt.Errorf("unexpected response %d: %s", resp.StatusCode(), string(resp.Body))
		return
	}
}

// PostMeme implements ClientInterface.
func (c Client) PostMeme(ctx context.Context, boardID models.BoardID, filename string, dsc map[string]string) (meme models.Meme, err error) {
	req := apiclient.PostMemeJSONRequestBody{BoardId: string(boardID), Filename: filename, Description: convertMapToAny(dsc)}
	resp, err := c.api.PostMemeWithResponse(ctx, req, c.middlewares()...)
	if err != nil {
		err = fmt.Errorf("can't request: %w", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		meme = convertMemeToModel(*resp.JSON200)
		return
	case 401:
		err = models.ErrUnauthorized
		return
	case 403:
		err = models.ErrForbidden
		return
	case 404:
		err = models.ErrMemeNotFound
		return
	default:
		err = fmt.Errorf("unexpected response %d: %s", resp.StatusCode(), string(resp.Body))
		return
	}
}

// PutMediaByID implements ClientInterface.
func (c Client) PutMediaByID(ctx context.Context, media models.Media) (err error) {
	contentType := http.DetectContentType(media.Body)
	resp, err := c.api.PutMediaByIDWithBodyWithResponse(ctx, apiclient.MediaId(media.ID), contentType, bytes.NewBuffer(media.Body), c.middlewares()...)
	if err != nil {
		err = fmt.Errorf("can't request: %w", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		return
	case 401:
		err = models.ErrUnauthorized
		return
	case 403:
		err = models.ErrForbidden
		return
	case 404:
		err = models.ErrMediaNotFound
		return
	default:
		err = fmt.Errorf("unexpected response %d: %s", resp.StatusCode(), string(resp.Body))
		return
	}
}

// SearchByBoardID implements ClientInterface.
func (c Client) SearchMemes(ctx context.Context, page int, pageSize int, general string) (memes []models.ScoredMeme, err error) {
	req := &apiclient.SearchMemesParams{Page: &page, PageSize: &pageSize, General: &general}
	resp, err := c.api.SearchMemesWithResponse(ctx, req, c.middlewares()...)
	if err != nil {
		err = fmt.Errorf("can't request: %w", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		return
	case 401:
		err = models.ErrUnauthorized
		return
	case 403:
		err = models.ErrForbidden
		return
	default:
		err = fmt.Errorf("unexpected response %d: %s", resp.StatusCode(), string(resp.Body))
		return
	}
}

// SubscribeByBoardID implements ClientInterface.
func (c Client) SubscribeByBoardID(ctx context.Context, boardID models.BoardID) (err error) {
	resp, err := c.api.SubscribeByBoardIDWithResponse(ctx, apiclient.BoardId(boardID), c.middlewares()...)
	if err != nil {
		err = fmt.Errorf("can't request: %w", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		return
	case 401:
		err = models.ErrUnauthorized
		return
	case 403:
		err = models.ErrForbidden
		return
	default:
		err = fmt.Errorf("unexpected response %d: %s", resp.StatusCode(), string(resp.Body))
		return
	}
}

// UnsubscribeByBoardID implements ClientInterface.
func (c Client) UnsubscribeByBoardID(ctx context.Context, boardID models.BoardID) (err error) {
	resp, err := c.api.UnsubscribeByBoardIDWithResponse(ctx, apiclient.BoardId(boardID), c.middlewares()...)
	if err != nil {
		err = fmt.Errorf("can't request: %w", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		return
	case 401:
		err = models.ErrUnauthorized
		return
	case 403:
		err = models.ErrForbidden
		return
	case 404:
		err = models.ErrSubNotFound
		return
	default:
		err = fmt.Errorf("unexpected response %d: %s", resp.StatusCode(), string(resp.Body))
		return
	}
}

// UpdateBoardByID implements ClientInterface.
func (c Client) UpdateBoardByID(ctx context.Context, boardID models.BoardID, name *string, owner *models.UserID) (board models.Board, err error) {
	req := apiclient.UpdateBoardByIDJSONRequestBody{Name: name, Owner: (*string)(owner)}
	resp, err := c.api.UpdateBoardByIDWithResponse(ctx, apiclient.BoardId(boardID), req, c.middlewares()...)
	if err != nil {
		err = fmt.Errorf("can't request: %w", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		board = convertBoardToModel(*resp.JSON200)
		return
	case 401:
		err = models.ErrUnauthorized
		return
	case 403:
		err = models.ErrForbidden
		return
	case 400:
		err = parseApiError(*resp.JSON400)
		return
	case 404:
		err = models.ErrBoardNotFound
		return
	default:
		err = fmt.Errorf("unexpected response %d: %s", resp.StatusCode(), string(resp.Body))
		return
	}
}

// UpdateMemeByID implements ClientInterface.
func (c Client) UpdateMemeByID(ctx context.Context, memeID models.MemeID, boardID *models.BoardID, filename *string, dsc *map[string]string) (meme models.Meme, err error) {
	req := apiclient.UpdateMemeByIDJSONRequestBody{BoardId: (*string)(boardID), Filename: filename}
	if dsc != nil {
		req.Description = ptr(convertMapToAny(*dsc))
	}
	resp, err := c.api.UpdateMemeByIDWithResponse(ctx, apiclient.MemeId(memeID), req, c.middlewares()...)
	if err != nil {
		err = fmt.Errorf("can't request: %w", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		meme = convertMemeToModel(*resp.JSON200)
		return
	case 401:
		err = models.ErrUnauthorized
		return
	case 403:
		err = models.ErrForbidden
		return
	case 404:
		err = models.ErrMemeNotFound
		return
	default:
		err = fmt.Errorf("unexpected response %d: %s", resp.StatusCode(), string(resp.Body))
		return
	}
}
