package statemachine

import (
	"api-client/pkg/models"
	"errors"
	"fmt"
)

func doLogin(r RequestContext, login, password string) error {
	ctx := r.Ctx

	token, err := r.ApiClient.AuthLogin(ctx, login, password)
	if err != nil {
		return fmt.Errorf("can't login: %w", err)
	}
	r.ApiClient.SetToken(token)
	err = doWhoami(r)
	if err != nil {
		return fmt.Errorf("can't do whoami: %w", err)
	}
	return nil
}

func doWhoami(r RequestContext) error {
	ctx := r.Ctx

	user, err := r.ApiClient.AuthWhoami(ctx)
	if err != nil {
		return fmt.Errorf("can't whoami: %w", err)
	}

	msg := fmt.Sprintf(`Logged in as %s (<code>%s</code>)`, user.Login, user.ID)
	_, err = r.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}
	return nil
}

func doGetBoard(r RequestContext) error {
	ctx := r.Ctx

	id := models.BoardID(r.UserInfo.activeBoard)
	b, err := r.ApiClient.GetBoardByID(ctx, id)
	if err != nil {
		return fmt.Errorf("can't get board by id: %w", err)
	}

	msg := fmt.Sprintf(`Active board is %s (<code>%s</code>)`, b.Name, b.ID)
	_, err = r.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}
	return nil
}

func doSetBoard(r RequestContext, id models.BoardID) error {
	r.UserInfo.activeBoard = id

	err := doGetBoard(r)
	if err != nil {
		return fmt.Errorf("can't doGetBoard: %w", err)
	}
	return nil
}

func doSubscribe(r RequestContext, id models.BoardID) error {
	ctx := r.Ctx
	err := r.ApiClient.SubscribeByBoardID(ctx, id)
	if err != nil {
		return fmt.Errorf("can't subscribe: %w", err)
	}

	_, err = r.SendMessage("Success")
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}
	return nil
}
func doUnsubscribe(r RequestContext, id models.BoardID) error {
	ctx := r.Ctx
	err := r.ApiClient.UnsubscribeByBoardID(ctx, id)
	if err != nil {
		return fmt.Errorf("can't subscribe: %w", err)
	}

	_, err = r.SendMessage("Success")
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}
	return nil
}

func sendError(r RequestContext, err error) {
	switch {
	case errors.Is(err, models.ErrBoardNotFound),
		errors.Is(err, models.ErrForbidden),
		errors.Is(err, models.ErrLoginExists),
		errors.Is(err, models.ErrMediaIsRequired),
		errors.Is(err, models.ErrMediaNotFound),
		errors.Is(err, models.ErrMemeNotFound),
		errors.Is(err, models.ErrSubNotFound),
		errors.Is(err, models.ErrUnauthorized),
		errors.Is(err, models.ErrUserNotFound),
		errors.Is(err, models.ErrInvalidInput{}):
		r.SendMessage(unwrap(err).Error())

	}
}

func unwrap(e error) error {
	type u interface {
		Unwrap() error
	}

	for {
		ne, ok := e.(u)
		if !ok {
			break
		}
		e = ne.Unwrap()
	}
	return e
}
