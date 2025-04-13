//go:generate go tool oapi-codegen -config cfg.yaml ../../../openapi.yaml
package apiserver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"memesearch/internal/api"
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
	switch {
	case errors.Is(err, api.ErrMemeNotFound):
		return nil, ErrMemeNotFound
	case err != nil:
		return nil, fmt.Errorf("can't delete meme: %w", err)
	}
	return DeleteMemeByID204Response{}, nil
}

// GetMediaByID implements StrictServerInterface.
func (s ServerImpl) GetMediaByID(ctx context.Context, request GetMediaByIDRequestObject) (GetMediaByIDResponseObject, error) {
	media, err := s.api.GetMedia(ctx, models.MediaID(request.Id))
	switch {
	case errors.Is(err, api.ErrMediaNotFound):
		return nil, ErrMediaNotFound
	case err != nil:
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
	switch {
	case errors.Is(err, api.ErrMemeNotFound):
		return nil, ErrMemeNotFound
	case err != nil:
		return nil, fmt.Errorf("can't get meme")
	}
	return GetMemeByID200JSONResponse(castMemesFromModel(meme)), nil
}

// ListMemes implements StrictServerInterface.
func (s ServerImpl) ListMemes(ctx context.Context, request ListMemesRequestObject) (ListMemesResponseObject, error) {
	page := 1
	pageSize := 10

	if request.Params.Page != nil {
		page = *request.Params.Page
	}
	if request.Params.PageSize != nil {
		pageSize = *request.Params.PageSize
	}

	if page < 1 || pageSize < 1 || pageSize > 100 {
		return nil, ErrInvalidPagination
	}

	//TODO sortBy

	return ListMemes200JSONResponse{
		Items:    []Meme{},
		Page:     page,
		PageSize: pageSize,
		Total:    -1,
	}, nil
}

// PostMeme implements StrictServerInterface.
func (s ServerImpl) PostMeme(ctx context.Context, request PostMemeRequestObject) (PostMemeResponseObject, error) {
	if request.Body == nil {
		return nil, &InvalidParamFormatError{ParamName: "body", Err: fmt.Errorf("empty body")}
	}

	meme := models.Meme{BoardID: models.BoardID(request.Body.BoardId)}
	if request.Body.Description != nil {
		dscs := map[string]string{}
		for k, va := range *request.Body.Description {
			v, ok := va.(string)
			if !ok {
				return nil, &InvalidParamFormatError{ParamName: "description", Err: fmt.Errorf("description must be map[string]string")}
			}
			dscs[k] = v
		}
		meme.Descriptions = dscs
	}
	if request.Body.Filename != nil {
		meme.Filename = *request.Body.Filename
	}

	id, err := s.api.CreateMeme(ctx, meme)
	if err != nil {
		return nil, fmt.Errorf("can't create meme: %w", err)
	}

	return PostMeme201JSONResponse{Id: string(id)}, nil
}

// PutMediaByID implements StrictServerInterface.
func (s ServerImpl) PutMediaByID(ctx context.Context, request PutMediaByIDRequestObject) (PutMediaByIDResponseObject, error) {
	form, err := request.Body.ReadForm(18 * 1024 * 1024) // 18MB limit
	if err != nil {
		if err == multipart.ErrMessageTooLarge {
			return nil, ErrTooLarge
		}
		return nil, fmt.Errorf("can't read form: %w", err)
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
		return nil, ErrUnsupportedMediaType
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
	meme := models.Meme{ID: models.MemeID(request.Id)}
	if request.Body.Description != nil {
		dscs := map[string]string{}
		for k, va := range *request.Body.Description {
			v, ok := va.(string)
			if !ok {
				return nil, &InvalidParamFormatError{ParamName: "description", Err: fmt.Errorf("description must be map[string]string")}
			}
			dscs[k] = v
		}
		meme.Descriptions = dscs
	}
	if request.Body.Filename != nil {
		meme.Filename = *request.Body.Filename
	}
	if request.Body.BoardId != nil {
		meme.BoardID = models.BoardID(*request.Body.BoardId)
	}

	err := s.api.UpdateMeme(ctx, meme)
	if err != nil {
		if errors.Is(err, api.ErrMemeNotFound) {
			return nil, ErrMemeNotFound
		}
		return nil, fmt.Errorf("can't update meme: %w", err)
	}
	meme, err = s.api.GetMemeByID(ctx, models.MemeID(request.Id))
	if err != nil {
		return nil, fmt.Errorf("can't get meme: %w", err)
	}
	return UpdateMemeByID200JSONResponse(castMemesFromModel(meme)), nil
}

func castMemesFromModel(meme models.Meme) Meme {
	dsc := map[string]any{}
	for k, v := range meme.Descriptions {
		dsc[k] = v
	}
	return Meme{
		Id:          string(meme.ID),
		BoardId:     string(meme.BoardID),
		Filename:    &meme.Filename, //TODO emptyfilename?
		CreatedAt:   meme.CreatedAt,
		UpdatedAt:   meme.UpdatedAt,
		Description: &dsc,
	}
}
