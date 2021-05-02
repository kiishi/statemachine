package statemachine

type TransitionRule struct {
	EventName        string
	CurrentState     string
	DestinationState string
	OnTransition     func(rule *TransitionRule, newState State)
}
