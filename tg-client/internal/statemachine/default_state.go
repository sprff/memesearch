package statemachine

var _ State = &DefaultState{}

type DefaultState struct {
}

func (d *DefaultState) Process(r RequestContext) (State, error) {
	chat := *(r.Event.FromChat())
	r.Bot.SendMessage(r.Ctx, chat.ID, "Default")
	return &CentralState{}, nil
}
