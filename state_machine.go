package statemachine

import (
	"errors"
)

type StateMachine struct {
	CurrentState          State
	StateMap              map[string]State
	TransitionRules       map[string]*TransitionRule
}

type Config struct {
	States      []State
	Transitions []*TransitionRule
}

func NewMachine(config *Config, initialStateIndex int) *StateMachine {
	var defaultCurrentState State
	if config.States != nil && len(config.States) > initialStateIndex{
		defaultCurrentState = config.States[initialStateIndex]
	}

	newMachine:= &StateMachine{
		CurrentState:          defaultCurrentState,
		StateMap:              make(map[string]State),
		TransitionRules:       make(map[string]*TransitionRule),
	}

	if config.States != nil{
		for _ , val := range config.States{
			newMachine.AddState(val)
		}
	}

	if config.Transitions != nil{
		for _ , val := range config.Transitions{
			newMachine.AddTransition(val)
		}
	}

	return newMachine
}

func (s *StateMachine) AddTransition(transitionRule *TransitionRule) {
	s.TransitionRules[transitionRule.EventName] = transitionRule
}

func (s *StateMachine) AddState(state State) error {
	// if the current state is nil, the first state added becomes the default state
	if s.CurrentState == nil{
		s.CurrentState = state
	}
	s.StateMap[state.GetIdentifier()] = state
	return nil
}

func (s *StateMachine) EmitSequence(eventNames ...string) error{
	for _, event := range eventNames{
		err := s.Emit(event)
		if err != nil{
			return errors.New("Invalid Transaction sequence ==>"+ err.Error())
		}
	}
	return nil
}

func (s *StateMachine) SetState(stateId string) error{
	if state , ok := s.StateMap[stateId]; ok{
		s.CurrentState = state
		return nil
	}
	return errors.New("State doesn't exist")
}

func (s *StateMachine) Emit(eventName string) error {
	if transition, ok := s.TransitionRules[eventName]; ok {
		if transition.CurrentState == s.CurrentState.GetIdentifier() {
			//	check if destination state exist
			if val, ok := s.StateMap[transition.DestinationState]; ok {
				if transition.OnTransition != nil{
					transition.OnTransition(transition, val, s.CurrentState)
				}
				s.CurrentState = val
				return nil
			}
			return errors.New("Destination State not registered")
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
