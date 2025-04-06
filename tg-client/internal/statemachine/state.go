package statemachine

type State interface {
	Process(r RequestContext) (State, error)
}
