package lib

type TransitionFunction = func(state interface{})
type TransitionCallbackFunction = func(state interface{}) interface{}
type Transition struct {
	Name     string
	From     string
	To       string
	NewState interface{}
	Action   TransitionFunction
}