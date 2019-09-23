package lib

import (
	"errors"
	"fmt"
)

type StateMachine struct {
	Initial_State         interface{}
	transitionMap         map[string]*Transition
	transitionFunctionMap map[string]TransitionFunction
}

func NewMachine(initial_state string) *StateMachine {
	return &StateMachine{
		Initial_State:         initial_state,
		transitionFunctionMap: make(map[string]TransitionFunction),
		transitionMap:         make(map[string]*Transition),
	}
}

func (this *StateMachine) SetState(value interface{}) {
	this.Initial_State = value
}

func (this *StateMachine) AddTransition(transiton *Transition) {
	this.transitionMap[transiton.Name] = transiton
	this.transitionFunctionMap[transiton.Name] = transiton.Action
}

func (this *StateMachine) RunTransition(ActionName string, callback TransitionFunction) error {
	if name, ok := this.transitionFunctionMap[ActionName]; ok {
		if this.transitionMap[ActionName].From != this.Initial_State {
			fmt.Println("Cannot transit")
			// remain in this state
			return nil
		} else {
			this.SetState(this.transitionMap[ActionName].To)
			name(this.Initial_State)
			callback(this.Initial_State)
			return nil
		}
	} else {
		return errors.New("Transitions Not Found")
	}
}

func (this *StateMachine) RunTransitionSequence(transitionSequence []string) {
	for _, value := range transitionSequence {
		this.RunTransition(value, func(state interface{}) {})
	}
}
