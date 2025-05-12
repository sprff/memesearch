package statemachine

import (
	"api-client/pkg/models"
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

func doRegister(r RequestContext, login, password string) error {
	ctx := r.Ctx

	_, err := r.ApiClient.AuthRegister(ctx, login, password)
	if err != nil {
		return fmt.Errorf("can't register: %w", err)
	}
	err = doLogin(r, login, password)
	if err != nil {
		return fmt.Errorf("can't do whoami: %w", err)
	}
	return nil
}
func doLogin(r RequestContext, login, password string) error {
	ctx := r.Ctx

	token, err := r.ApiClient.AuthLogin(ctx, login, password)
	if err != nil {
		return fmt.Errorf("can't login: %w", err)
	}
	r.UserInfo.Token = token
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

	id := models.BoardID(r.UserInfo.ActiveBoard)
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
	prev := r.UserInfo.ActiveBoard
	r.UserInfo.ActiveBoard = id

	err := doGetBoard(r)
	if err != nil {
		r.UserInfo.ActiveBoard = prev
		return fmt.Errorf("can't doGetBoard: %w", err)
	}
	return nil
}

func doCreateBoard(r RequestContext, name string) error {
	ctx := r.Ctx
	b, err := r.ApiClient.PostBoard(ctx, name)
	if err != nil {
		return fmt.Errorf("can't post board: %w", err)
	}

	msg := fmt.Sprintf("New board: %s (<code>%s</code>)", b.Name, b.ID)
	_, err = r.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}
	return nil
}
func doListBoards(r RequestContext) error {
	ctx := r.Ctx
	boards, err := r.ApiClient.ListBoards(ctx, 0, 100, "id")
	if err != nil {
		return fmt.Errorf("can't list boards: %w", err)
	}

	msg := strings.Builder{}
	for i, b := range boards {
		msg.WriteString(fmt.Sprintf("%d. %s (<code>%s</code>)\n", i+1, b.Name, b.ID))
	}
	_, err = r.SendMessage(msg.String())
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
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
		errors.Is(err, ErrBadCommandUsage),
		errors.Is(err, models.ErrInvalidInput{}):
		r.SendMessage(unwrap(err).Error())
		slog.InfoContext(r.Ctx, "Expected error", "err", unwrap(err).Error())
	default:
		id := r.ApiClient.GetID()
		r.SendMessage(fmt.Sprintf("unexpected error: %s", id))
		slog.ErrorContext(r.Ctx, "Error in processing", "err", err, "event", r.Event)

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

func help() string {
	return `MemeSearch - бот для поиска мемов по описанию
1) Поиск: осуществляется командой /search query либо inline запросом @MemeManiac query. Для поиска по видео используйте @MemeManiac !query
2) Бот учитывает аккаунт(сервиса MemeSearch, не телегерама) с которого приходят запросы и использует мемы доступные этому аккаунту.
3) Команды для работы с аккаунтом:
	/register login password - регистраиция
	/login login password - авторизация
	/whoami - узнать что за аккаунт сейчас активен
	/logout - выйти из аккаунта
4) Сервис поддерживает концепцию досок мемов:
	/getboard - узнать текущую активную доску
	/setboard id - указать новую активную доску
	/createboard name - Создать доску с именем name
	/listboards - Перечислить доступные доски
	/subscibe id - Подписаться на доску id чтобы иметь доступ к ее мемам
	/unsubscribe id - Отписаться от доски id
5) Для того чтобы создать мем, пришлите фото/виде с описанием. Данный мем будет создан на текущую активную доску
`
}
