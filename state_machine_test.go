package statemachine

import (
	"bytes"
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
		})
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
		})
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
		}) //setting state1 as inital state
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
		})
		sampleStateMachine.EmitSequence("move", "jump")

		if sampleStateMachine.GetCurrentState().GetIdentifier() != strconv.Itoa(STATE3) {
			t.Fail()
		}
	})

	t.Run("OnTransition is called when a new state is entered", func(t *testing.T) {
		buffer := new(bytes.Buffer)
		transition1 := TransitionRule{
			EventName:        "move",
			CurrentState:     strconv.Itoa(STATE1),
			DestinationState: strconv.Itoa(STATE2),
			OnTransition: func(rule *TransitionRule, newState, previousState State) {
				buffer.Write([]byte("x"))
			},
		}
		transition2 := TransitionRule{
			EventName:        "jump",
			CurrentState:     strconv.Itoa(STATE2),
			DestinationState: strconv.Itoa(STATE3),
			OnTransition: func(rule *TransitionRule, newState, previousState State) {
				buffer.Write([]byte("y"))
			},
		}
		sampleStateMachine := NewMachine(&Config{
			States: []State{&state1, &state2, &state3},
			Transitions: []*TransitionRule{
				&transition1,
				&transition2,
			},
		})
		sampleStateMachine.EmitSequence("move", "jump")
		if buffer.String() != "xy" {
			t.Fail()
		}
	})

	t.Run("Should panic for empty state list in config ", func(t *testing.T) {
		transition1 := TransitionRule{
			EventName:        "move",
			CurrentState:     strconv.Itoa(STATE1),
			DestinationState: strconv.Itoa(STATE2),
		}
		defer func() {
			if r := recover(); r == nil {
				t.Error("Panic for invalid state not thrown")
			}
		}()

		NewMachine(&Config{
			States: []State{},
			Transitions: []*TransitionRule{
				&transition1,
			},
		})
	})

	t.Run("should allow heterogenous forms of transitionRule parameters (current and destination state)", func(t *testing.T) {
		transition1 := TransitionRule{
			EventName:        "move",
			CurrentState:     strconv.Itoa(STATE1),
			DestinationState: &state2, // uses pointer to state
		}
		sampleStateMachine := NewMachine(&Config{
			States: []State{&state1, &state2, &state3},
			Transitions: []*TransitionRule{
				&transition1,
			},
		})
		sampleStateMachine.Emit("move")
		if sampleStateMachine.GetCurrentState().GetIdentifier() != "1" {
			t.Errorf("Current state at %s", sampleStateMachine.GetCurrentState().GetIdentifier())
		}
	})

}
