package apiserver

import (
	"fmt"
	"memesearch/internal/models"
	"regexp"
	"slices"
)

func (r UpdateMemeByIDRequestObject) GetParams() (
	id models.MemeID, dsc *map[string]string, filename *string, board *models.BoardID, err error) {
	id = models.MemeID(r.Id)
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
		//TODO validate if board exists
		board = ptr(models.BoardID(*u.BoardId))
	}

	return
}

const (
	DefaultPageSize = 20
	DefaultPage     = 1
	DefaultSortBy   = "id"
)

var (
	AllowedSortBy = []string{"id", "createdAt", "updatedAt"}
)

func (r SearchByBoardIDRequestObject) GetParams() (
	id models.BoardID, page, pageSize int, sortBy string, dsc map[string]string, err error) {
	page = DefaultPage
	pageSize = DefaultPageSize
	sortBy = DefaultSortBy

	//TODO validate board
	id = models.BoardID(r.Id)

	if r.Params.Page != nil {
		page = *r.Params.Page
	}
	if page < 1 {
		err = invalidInput("page", "must be page>=1")
		return
	}

	if r.Params.PageSize != nil {
		pageSize = *r.Params.PageSize
	}
	if pageSize < 1 || pageSize > 100 {
		err = invalidInput("pageSize", "must be 1 <= pageSize <= 100")
		return
	}

	if r.Params.SortBy != nil {
		sortBy = string(*r.Params.SortBy)
	}
	if !slices.Contains(AllowedSortBy, sortBy) {
		err = invalidInput("sortBy", "sortBy must be one of %v", AllowedSortBy)
		return
	}

	dsc = getDescriptionMap(r.Params)
	return
}

func (r ListMemesRequestObject) GetParams() (
	page, pageSize int, sortBy string, err error) {
	page = DefaultPage
	pageSize = DefaultPageSize
	sortBy = DefaultSortBy

	if r.Params.Page != nil {
		page = *r.Params.Page
	}
	if page < 1 {
		err = invalidInput("page", "must be page>=1")
		return
	}

	if r.Params.PageSize != nil {
		pageSize = *r.Params.PageSize
	}
	if pageSize < 1 || pageSize > 100 {
		err = invalidInput("pageSize", "must be 1 <= pageSize <= 100")
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

func getDescriptionMap(p SearchByBoardIDParams) map[string]string {
	m := map[string]string{}
	if p.General != nil {
		m["general"] = *p.General
	}
	return m
}

func (r PostMemeRequestObject) GetParams() (
	id models.BoardID, filename string, dsc map[string]string, err error) {
	if r.Body == nil {
		err = invalidInput("body", "not empty body is expected")
		return
	}

	//TODO validate boardID
	id = models.BoardID(r.Body.BoardId)
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
