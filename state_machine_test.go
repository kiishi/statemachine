package main

import (
	"fmt"
	"strconv"
	"testing"
)

const (
	STATE1 = iota
	STATE2
	STATE3
)

type Base struct {
	Id string
}

func (b *Base) GetIdentifier() string {
	return b.Id
}

type State1 struct {
	Base
}
type State2 struct {
	Base
}
type State3 struct {
	Base
}

func TestNewMachine(t *testing.T) {
	state1 := State1{Base{Id: fmt.Sprintf("%d", STATE1)}}
	state2 := State2{Base{Id: fmt.Sprintf("%d", STATE2)}}
	state3 := State3{Base{Id: fmt.Sprintf("%d", STATE3)}}

	t.Run("Should do nothing when transition is not set and return error", func(t *testing.T) {
		sampleStateMachine := NewMachine(&Config{
			States:      []State{&state1, &state2, &state3},
			Transitions: nil,
		},
			0) //setting state1 as inital state
		sampleStateMachine.Emit("Hello")
		if sampleStateMachine.GetCurrentState().GetIdentifier() != state1.GetIdentifier() {
			t.Errorf("State Do not match")
		}
	})

	t.Run("Should transitions to the appropriate state", func(t *testing.T) {
		transition1 := TransitionRule{
			EventName:        "move",
			CurrentState:     strconv.Itoa(STATE1),
			DestinationState: strconv.Itoa(STATE2),
		}
		sampleStateMachine := NewMachine(&Config{
			States: []State{&state1, &state2, &state3},
			Transitions: []*TransitionRule{
				&transition1,
			},
		},
			0) //setting state1 as inital state
		sampleStateMachine.Emit("move")
		if sampleStateMachine.GetCurrentState().GetIdentifier() != "1" {
			t.Errorf("Current state at %s", sampleStateMachine.GetCurrentState().GetIdentifier())
		}
	})

	t.Run("Should return error when emitting to a state that doesnt exist", func(t *testing.T) {
		transition1 := TransitionRule{
			EventName:        "move",
			CurrentState:     strconv.Itoa(STATE1),
			DestinationState: strconv.Itoa(STATE3),
		}
		sampleStateMachine := NewMachine(&Config{
			States: []State{&state1, &state2}, //omit state 3
			Transitions: []*TransitionRule{
				&transition1,
			},
		},
			0) //setting state1 as inital state
		error := sampleStateMachine.Emit("move")
		//fmt.Println(sampleStateMachine.CurrentState.GetIdentifier() == strconv.Itoa(STATE1))

		if error == nil {
			t.Fail()
		}
	})

	t.Run("Should run transitions in sequence", func(t *testing.T) {
		transition1 := TransitionRule{
			EventName:        "move",
			CurrentState:     strconv.Itoa(STATE1),
			DestinationState: strconv.Itoa(STATE2),
		}
		transition2 := TransitionRule{
			EventName:        "jump",
			CurrentState:     strconv.Itoa(STATE2),
			DestinationState: strconv.Itoa(STATE3),
		}
		sampleStateMachine := NewMachine(&Config{
			States: []State{&state1, &state2, &state3},
			Transitions: []*TransitionRule{
				&transition1,
				&transition2,
			},
		},
			0)
		sampleStateMachine.EmitSequence("move", "jump")

		if sampleStateMachine.GetCurrentState().GetIdentifier() != strconv.Itoa(STATE3){
			t.Fail()
		}
	})
}
