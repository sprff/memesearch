//go:generate go tool oapi-codegen -config cfg.yaml ../../../openapi.yaml
package apiserver

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"memesearch/internal/api"
	"memesearch/internal/apiserver/middleware"
	"memesearch/internal/models"
	"mime/multipart"
	"net/http"
)

var _ StrictServerInterface = ServerImpl{}

type ServerImpl struct {
	api *api.API
}

func NewServerImpl(api *api.API) ServerImpl {
	return ServerImpl{
		api: api,
	}
}

// About implements StrictServerInterface.
func (s ServerImpl) About(ctx context.Context, request AboutRequestObject) (AboutResponseObject, error) {
	return About200JSONResponse{
		ApiName:     "Meme Search API",
		Description: "API for managing internet memes and related media",
		Version:     "1.0.0",
	}, nil
}

// DeleteMemeByID implements StrictServerInterface.
func (s ServerImpl) DeleteMemeByID(ctx context.Context, request DeleteMemeByIDRequestObject) (DeleteMemeByIDResponseObject, error) {
	err := s.api.DeleteMeme(ctx, models.MemeID(request.Id))
	if err != nil {
		return nil, fmt.Errorf("can't delete meme: %w", err)
	}
	return DeleteMemeByID204Response{}, nil
}

// GetMediaByID implements StrictServerInterface.
func (s ServerImpl) GetMediaByID(ctx context.Context, request GetMediaByIDRequestObject) (GetMediaByIDResponseObject, error) {
	media, err := s.api.GetMedia(ctx, models.MediaID(request.Id))
	if err != nil {
		return nil, fmt.Errorf("can't get media: %w", err)
	}

	mime := http.DetectContentType(media.Body[:512])
	buf := bytes.NewBuffer(media.Body)
	clen := int64(len(media.Body))
	switch mime {
	case "image/jpeg":
		return GetMediaByID200ImagejpegResponse{Body: buf, ContentLength: clen}, nil
	case "image/jpg":
		return GetMediaByID200ImagejpgResponse{Body: buf, ContentLength: clen}, nil
	case "image/png":
		return GetMediaByID200ImagepngResponse{Body: buf, ContentLength: clen}, nil
	case "video/mp4":
		return GetMediaByID200Videomp4Response{Body: buf, ContentLength: clen}, nil
	default:
		return GetMediaByID200ApplicationoctetStreamResponse{Body: buf, ContentLength: clen}, nil
	}

}

// GetMemeByID implements StrictServerInterface.
func (s ServerImpl) GetMemeByID(ctx context.Context, request GetMemeByIDRequestObject) (GetMemeByIDResponseObject, error) {
	meme, err := s.api.GetMemeByID(ctx, models.MemeID(request.Id))
	if err != nil {
		return nil, fmt.Errorf("can't get meme: %w", err)
	}
	return GetMemeByID200JSONResponse(convertMemeToServer(meme)), nil
}

// ListMemes implements StrictServerInterface.
func (s ServerImpl) ListMemes(ctx context.Context, request ListMemesRequestObject) (ListMemesResponseObject, error) {
	page, pageSize, sortBy, err := request.GetParams()
	if err != nil {
		return nil, fmt.Errorf("can't get params: %w", err)
	}

	memes, err := s.api.ListMemes(ctx, (page-1)*pageSize, pageSize, sortBy)
	if err != nil {
		return nil, fmt.Errorf("can't list memes: %w", err)
	}

	conv := make([]Meme, 0, len(memes))
	for _, m := range memes {
		conv = append(conv, convertMemeToServer(m))
	}

	return ListMemes200JSONResponse{Items: conv}, nil
}

// PostMeme implements StrictServerInterface.
func (s ServerImpl) PostMeme(ctx context.Context, request PostMemeRequestObject) (PostMemeResponseObject, error) {
	board, filename, dsc, err := request.GetParams()
	if err != nil {
		return nil, fmt.Errorf("can't get params: %w", err)
	}

	meme := models.Meme{BoardID: board, Filename: filename, Description: dsc}
	id, err := s.api.CreateMeme(ctx, meme)
	if err != nil {
		return nil, fmt.Errorf("can't create meme: %w", err)
	}

	return PostMeme201JSONResponse{Id: string(id)}, nil
}

// PutMediaByID implements StrictServerInterface.
func (s ServerImpl) PutMediaByID(ctx context.Context, request PutMediaByIDRequestObject) (PutMediaByIDResponseObject, error) {
	// TODO prettyfiy put media
	form, err := request.Body.ReadForm(20 * 1024 * 1024)
	if err != nil {
		if err == multipart.ErrMessageTooLarge {
			return nil, invalidInput("body", "body is too big")
		}
		return nil, invalidInput("form", "%s", err.Error())
	}

	files, ok := form.File["media"]
	if !ok || len(files) == 0 {
		return nil, &InvalidParamFormatError{ParamName: "form data", Err: fmt.Errorf("no file provided")}
	}

	fileHeader := files[0]
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("can't open file: %w", err)
	}
	defer file.Close()

	const maxFileSize = 16 * 1024 * 1024
	if fileHeader.Size > maxFileSize {
		return nil, &InvalidParamFormatError{ParamName: "form file", Err: fmt.Errorf("file size exceed maximum size")}
	}

	// 3. Validate media type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("can't read buffer prefix: %w", err)
	}

	contentType := http.DetectContentType(buffer)
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"video/mp4":  true,
	}

	if !allowedTypes[contentType] {
		return nil, invalidInput("media", "bad media type")
	}

	// Reset file pointer after reading the header
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("can't seek file: %w", err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("can't read file: %w", err)
	}

	err = s.api.SetMedia(ctx, models.Media{ID: models.MediaID(request.Id), Body: data})
	if err != nil {
		return nil, fmt.Errorf("can't set media: %w", err)
	}

	return PutMediaByID200Response{}, nil
}

// UpdateMemeByID implements StrictServerInterface.
func (s ServerImpl) UpdateMemeByID(ctx context.Context, request UpdateMemeByIDRequestObject) (UpdateMemeByIDResponseObject, error) {
	id, dsc, filename, board, err := request.GetParams()
	if err != nil {
		return nil, fmt.Errorf("can't get params: %w", err)
	}

	meme, err := s.api.GetMemeByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("can't get meme: %w", err)
	}
	if dsc != nil {
		meme.Description = *dsc
	}
	if filename != nil {
		meme.Filename = *filename
	}
	if board != nil {
		meme.BoardID = *board
	}

	err = s.api.UpdateMeme(ctx, meme)
	if err != nil {
		return nil, fmt.Errorf("can't update meme: %w", err)
	}

	meme, err = s.api.GetMemeByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("can't get meme: %w", err)
	}

	return UpdateMemeByID200JSONResponse(convertMemeToServer(meme)), nil
}

// SearchByBoardID implements StrictServerInterface.
func (s ServerImpl) SearchByBoardID(ctx context.Context, request SearchByBoardIDRequestObject) (SearchByBoardIDResponseObject, error) {
	boardID, page, pageSize, sortBy, dsc, err := request.GetParams()
	if err != nil {
		return nil, fmt.Errorf("can't get params: %w", err)
	}

	memes, err := s.api.SearchMemeByBoardID(ctx, boardID, dsc, (page-1)*pageSize, pageSize, sortBy)
	if err != nil {
		return nil, fmt.Errorf("can't search: %w", err)
	}

	conv := make([]Meme, 0, len(memes))
	for _, m := range memes {
		conv = append(conv, convertMemeToServer(m))
	}

	return SearchByBoardID200JSONResponse{Items: conv}, nil

}

// AuthLogin implements StrictServerInterface.
func (s ServerImpl) AuthLogin(ctx context.Context, request AuthLoginRequestObject) (AuthLoginResponseObject, error) {
	login, password, err := request.GetParams()
	if err != nil {
		return nil, fmt.Errorf("can't get params: %w", err)
	}

	token, err := s.api.LoginUser(ctx, login, password)
	if err != nil {
		return nil, fmt.Errorf("can't login: %w", err)
	}

	return AuthLogin200JSONResponse{Token: token}, nil
}

// AuthRegister implements StrictServerInterface.
func (s ServerImpl) AuthRegister(ctx context.Context, request AuthRegisterRequestObject) (AuthRegisterResponseObject, error) {
	login, password, err := request.GetParams()
	if err != nil {
		return nil, fmt.Errorf("can't get params: %w", err)
	}

	id, err := s.api.PostUser(ctx, login, password)
	if err != nil {
		return nil, fmt.Errorf("can't login: %w", err)
	}

	return AuthRegister200JSONResponse{Id: string(id)}, nil
}

// AuthWhoami implements StrictServerInterface.
func (s ServerImpl) AuthWhoami(ctx context.Context, request AuthWhoamiRequestObject) (AuthWhoamiResponseObject, error) {
	id := middleware.GetAuthUserID(ctx)
	if id == "" {
		return AuthWhoami401Response{}, nil
	}
	user, err := s.api.GetUser(ctx, models.UserID(id))
	if err != nil {
		if err == api.ErrUserNotFound {
			return AuthWhoami401Response{}, nil
		}
		return nil, fmt.Errorf("can't get user: %w", err)
	}
	return AuthWhoami200JSONResponse{Id: string(user.ID), Login: user.Login}, nil
}

// GetUserByID implements StrictServerInterface.
func (s ServerImpl) GetUserByID(ctx context.Context, request GetUserByIDRequestObject) (GetUserByIDResponseObject, error) {
	id := request.Id
	user, err := s.api.GetUser(ctx, models.UserID(id))
	if err != nil {
		return nil, fmt.Errorf("can't get user: %w", err)
	}
	return GetUserByID200JSONResponse{Id: string(user.ID), Login: user.Login}, nil
}
