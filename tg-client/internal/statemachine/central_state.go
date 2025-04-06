package statemachine

import "log/slog"

var _ State = &CentralState{}

type CentralState struct {
}

func (d *CentralState) Process(r RequestContext) (State, error) {
	switch {
	case isSearchRequest(r):
		doSearchRequest(r)
		return &CentralState{}, nil
	case isInlineSearchRequest(r):
		doInlineSearchRequest(r)
		return &CentralState{}, nil
	case isAddPhoto(r):
		doAddPhoto(r)
		return &CentralState{}, nil
	case isAddVideo(r):
		doAddVideo(r)
		return &CentralState{}, nil
	default:
		return &CentralState{}, nil
	}
}

func isSearchRequest(r RequestContext) bool {
	if r.Event == nil || r.Event.Message == nil {
		return false
	}
	msg := r.Event.Message
	if len(msg.Photo) != 0 ||
		msg.Video != nil ||
		msg.Audio != nil ||
		msg.Document != nil ||
		msg.Voice != nil {
		return false
	}

	return msg.Text != ""
}

func doSearchRequest(r RequestContext) {
	slog.InfoContext(r.Ctx, "doSearchRequest")
	meme, err := r.ApiClient.GetMemeByID(r.Ctx, "0195d17ca1f47a01a4cc837f853da8f3")
	slog.InfoContext(r.Ctx, "meme request", "meme", meme, "err", err)
	r.Bot.SendMessage(r.MustChat(), "doSearchRequest")
}

func isInlineSearchRequest(r RequestContext) bool {
	return false
}

func doInlineSearchRequest(r RequestContext) {
	slog.InfoContext(r.Ctx, "doInlineSearchRequest")
}

func isAddPhoto(r RequestContext) bool {
	if r.Event == nil || r.Event.Message == nil {
		return false
	}
	msg := r.Event.Message
	if len(msg.Photo) == 0 ||
		msg.Video != nil ||
		msg.Audio != nil ||
		msg.Document != nil ||
		msg.Voice != nil {
		return false
	}

	return true
}

func doAddPhoto(r RequestContext) {
	slog.InfoContext(r.Ctx, "doAddPhoto")
	r.Bot.SendMessage(r.MustChat(), "doAddPhoto")
}

func isAddVideo(r RequestContext) bool {
	if r.Event == nil || r.Event.Message == nil {
		return false
	}
	msg := r.Event.Message
	if len(msg.Photo) != 0 ||
		msg.Video == nil ||
		msg.Audio != nil ||
		msg.Document != nil ||
		msg.Voice != nil {
		return false
	}

	return true
}

func doAddVideo(r RequestContext) {
	slog.InfoContext(r.Ctx, "doAddVideo")
	r.Bot.SendMessage(r.MustChat(), "doAddVideo")
}
