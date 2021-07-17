package statemachine

type TransitionRule struct {
	EventName        string
	CurrentState     interface{}
	DestinationState interface{}
	OnTransition     func(rule *TransitionRule, newState, previousState State)
}
