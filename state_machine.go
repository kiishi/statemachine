package statemachine

import (
	"errors"
	"fmt"
)

type StateMachine struct {
	CurrentState    State
	StateMap        map[string]State
	TransitionRules map[string]*TransitionRule
}

type Config struct {
	States      []State
	Transitions []*TransitionRule
}

func validateConfig(config *Config) {
	//	check state length
	if len(config.States) == 0 {
		panic(errors.New("Cannot accept Empty array for states"))
	}
	allStatesMap := map[string]int{}

	for _, val := range config.States {
		identifier := val.GetIdentifier()
		if _, ok := allStatesMap[identifier]; !ok {
			allStatesMap[identifier] = 1
			continue
		}
		panic(errors.New(fmt.Sprintf("Duplicate state identifier %s found ", identifier)))
	}

	/*	TODO: add check for transactionRule, should check if the currentState and
		destinationState are strictly strings or pointers to structs that implement the State interface
	*/

}

func NewMachine(config ...*Config) *StateMachine {
	//throw an error is state array is empty
	if len(config) == 0 {
		return &StateMachine{
			StateMap:        make(map[string]State),
			TransitionRules: make(map[string]*TransitionRule),
		}
	}
	validateConfig(config[0])
	newMachine := &StateMachine{
		CurrentState:    config[0].States[0],
		StateMap:        make(map[string]State),
		TransitionRules: make(map[string]*TransitionRule),
	}

	if config[0].States != nil {
		for _, val := range config[0].States {
			newMachine.AddState(val)
		}
	}

	if config[0].Transitions != nil {
		for _, val := range config[0].Transitions {
			newMachine.AddTransition(val)
		}
	}

	return newMachine
}

func (s *StateMachine) AddTransition(transitionRule *TransitionRule) *StateMachine {
	s.TransitionRules[transitionRule.EventName] = transitionRule
	return s
}

func (s *StateMachine) AddState(state State) *StateMachine {
	// if the current state is nil, the first state added becomes the default state
	if s.CurrentState == nil {
		s.CurrentState = state
	}
	s.StateMap[state.GetIdentifier()] = state
	return s
}

func (s *StateMachine) EmitInSequence(eventNames ...string) error {
	for _, event := range eventNames {
		err := s.Emit(event)
		if err != nil {
			return errors.New("Invalid Transaction sequence ==>" + err.Error())
		}
	}
	return nil
}

func (s *StateMachine) SetState(stateId string) error {
	if state, ok := s.StateMap[stateId]; ok {
		s.CurrentState = state
		return nil
	}
	return errors.New("State doesn't exist")
}

func (s *StateMachine) Emit(eventName string) error {
	if transition, ok := s.TransitionRules[eventName]; ok {
		var stringCurrentState string
		var stringDestinationState string

		/*
			cleanup transition rule by recognizing which of the state interfaces are
			regular strings or state structs
		*/
		if val, ok := transition.CurrentState.(string); ok {
			stringCurrentState = val
		} else {
			currentState, _ := transition.CurrentState.(State)
			stringCurrentState = currentState.GetIdentifier()
		}

		if val, ok := transition.DestinationState.(string); ok {
			stringDestinationState = val
		} else {
			destinationState, _ := transition.DestinationState.(State)
			stringDestinationState = destinationState.GetIdentifier()
		}

		// check if the current state is the expected state on the transition rule
		if transition.CurrentState == stringCurrentState {
			//	check if destination state exist
			if val, ok := s.StateMap[stringDestinationState]; ok {
				if transition.OnTransition != nil {
					transition.OnTransition(transition, val, s.CurrentState)
				}
				s.CurrentState = val
				return nil
			}
			return errors.New("Destination State does not exist")
		}
		//	if the current state is not what the transition expects
		return nil
	}
	//if there is no transition remain in the same state
	return nil
}

func (s *StateMachine) GetCurrentState() State {
	return s.CurrentState
}
