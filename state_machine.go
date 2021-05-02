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
	newMachine:= &StateMachine{
		CurrentState:          config.States[initialStateIndex],
		StateMap:              make(map[string]State),
		TransitionRules:       make(map[string]*TransitionRule),
	}
	for _ , val := range config.States{
		newMachine.AddState(val)
	}
	for _ , val := range config.Transitions{
		newMachine.AddTransition(val)
	}
	return newMachine
}

func (s *StateMachine) AddTransition(transitionRule *TransitionRule) {
	s.TransitionRules[transitionRule.EventName] = transitionRule
}

func (s *StateMachine) AddState(state State) error {
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

func (s *StateMachine) Emit(eventName string) error {
	if transition, ok := s.TransitionRules[eventName]; ok {
		if transition.CurrentState == s.CurrentState.GetIdentifier() {
			//	check if destination state exist
			if val, ok := s.StateMap[transition.DestinationState]; ok {
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
