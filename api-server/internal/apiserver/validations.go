package apiserver

import (
	"fmt"
	"memesearch/internal/models"
	"regexp"
	"slices"
)

func (r UpdateMemeByIDRequestObject) GetParams() (
	id models.MemeID, board *models.BoardID, filename *string, dsc *map[string]string, err error) {
	id = models.MemeID(r.MemeID)
	u := r.Body
	if u == nil {
		err = invalidInput("body", "not empty body is expected")
		return
	}

	if u.Description != nil {
		var dscs map[string]string
		dscs, err = convertMapToString(*u.Description)
		if err != nil {
			err = invalidInput("description", "description must be map[string]string")
			return
		}
		dsc = &dscs
	}
	if u.Filename != nil {
		filename = u.Filename
	}
	if u.BoardId != nil {
		board = ptr(models.BoardID(*u.BoardId))
	}

	return
}

const (
	DefaultOffset = 0
	DefaultLimit  = 20
	DefaultSortBy = "id"
)

var (
	AllowedSortBy = []string{"id", "createdAt", "updatedAt"}
)

func (r SearchMemesRequestObject) GetParams() (
	offset, limit int, dsc map[string]string, err error) {
	offset = DefaultOffset
	limit = DefaultLimit

	if r.Params.Offset != nil {
		offset = *r.Params.Offset
	}
	if offset < 0 {
		err = invalidInput("offset", "must be offset>=0")
		return
	}

	if r.Params.Limit != nil {
		limit = *r.Params.Limit
	}
	if limit < 1 || limit > 100 {
		err = invalidInput("limit", "must be 1 <= limit <= 100")
		return
	}

	dsc = getDescriptionMap(r.Params)
	return
}

func (r ListMemesRequestObject) GetParams() (
	offset, limit int, sortBy string, err error) {
	offset = DefaultOffset
	limit = DefaultLimit
	sortBy = DefaultSortBy

	if r.Params.Offset != nil {
		offset = *r.Params.Offset
	}
	if offset < 0 {
		err = invalidInput("offset", "must be offset>=0")
		return
	}

	if r.Params.Limit != nil {
		limit = *r.Params.Limit
	}
	if limit < 1 || limit > 100 {
		err = invalidInput("limit", "must be 1 <= limit <= 100")
		return
	}

	if r.Params.SortBy != nil {
		sortBy = string(*r.Params.SortBy)
	}
	if !slices.Contains(AllowedSortBy, sortBy) {
		err = invalidInput("sortBy", "sortBy must be one of %v", AllowedSortBy)
		return
	}
	return
}

func getDescriptionMap(p SearchMemesParams) map[string]string {
	m := map[string]string{}
	if p.General != nil {
		m["general"] = *p.General
	}
	return m
}

func (r PostMemeRequestObject) GetParams() (
	board models.BoardID, filename string, dsc map[string]string, err error) {
	if r.Body == nil {
		err = invalidInput("body", "not empty body is expected")
		return
	}

	board = models.BoardID(r.Body.BoardId)
	filename = r.Body.Filename
	dsc, err = convertMapToString(r.Body.Description)
	if err != nil {
		err = invalidInput("description", "description must be map[string]string")
		return
	}
	return
}

func validateLogin(login string) error {
	if len(login) < 3 || len(login) > 30 {
		return fmt.Errorf("login length must be in [3;30]")
	}
	format := `^[a-zA-Z0-9]*$`
	if ok := regexp.MustCompile(format).MatchString(login); !ok {
		return fmt.Errorf("login doesn't satisfy format: %s", format)
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password length must be at least 8")
	}

	return nil
}

func (r AuthLoginRequestObject) GetParams() (
	login, password string, err error) {
	login = r.Body.Login
	err = validateLogin(login)
	if err != nil {
		err = invalidInput("login", "%s", err.Error())
		return
	}

	password = r.Body.Password
	return
}

func (r AuthRegisterRequestObject) GetParams() (
	login, password string, err error) {
	login = r.Body.Login
	err = validateLogin(login)
	if err != nil {
		err = invalidInput("login", "%s", err.Error())
		return
	}

	password = r.Body.Password
	err = validatePassword(password)
	if err != nil {
		err = invalidInput("login", "%s", err.Error())
		return
	}
	return
}

func (r PostBoardRequestObject) GetParams() (
	name string, err error) {
	name = r.Body.Name
	if len(name) < 3 || 30 < len(name) {
		err = invalidInput("name", "name length should be [3;30]")
		return
	}
	return
}

func (r UpdateBoardByIDRequestObject) GetParams() (
	id models.BoardID, name *string, owner *models.UserID, err error) {
	id = models.BoardID(r.BoardID)
	name = r.Body.Name
	if name != nil && (len(*name) < 3 || 30 < len(*name)) {
		err = invalidInput("name", "name length should be [3;30]")
		return
	}
	if r.Body.Owner != nil {
		owner = ptr(models.UserID(*r.Body.Owner))
	}
	return
}
func (r ListBoardsRequestObject) GetParams() (
	offset, limit int, sortBy string, err error) {
	offset = DefaultOffset
	limit = DefaultLimit
	sortBy = DefaultSortBy

	if r.Params.Offset != nil {
		offset = *r.Params.Offset
	}
	if offset < 0 {
		err = invalidInput("offset", "must be offset>=0")
		return
	}

	if r.Params.Limit != nil {
		limit = *r.Params.Limit
	}
	if limit < 1 || limit > 100 {
		err = invalidInput("limit", "must be 1 <= limit <= 100")
		return
	}

	if r.Params.SortBy != nil {
		sortBy = string(*r.Params.SortBy)
	}
	if !slices.Contains(AllowedSortBy, sortBy) {
		err = invalidInput("sortBy", "sortBy must be one of %v", AllowedSortBy)
		return
	}
	return
}
