## State Machine ![](https://github.com/kiishi/statemachine/workflows/CI/badge.svg)
Cool-ish state machine💅🏾

### Sample Usage

```go
package main

import (
	"fmt"
	machine "github.com/kiishi/statemachine"
)

type FoodState struct {
	name string
}

func (m *FoodState) GetIdentifier() string {
	return m.name; // name must be unique
}

func main() {
	raw := FoodState{"raw"}
	cooked := FoodState{"cooked"}

	statemachine := machine.NewMachine(
		&machine.Config{
			States: []machine.State{
				&raw, // first state is initial
				&cooked,
			},
			Transitions: []*machine.TransitionRule{
				&machine.TransitionRule{
					EventName:        "cook",
					CurrentState:     &raw,
					DestinationState: &cooked,
					OnTransition: func(
						rule *machine.TransitionRule,
						newState,
						previousState machine.State,
					) {
						fmt.Println("Food is ready")
					},
				},
			},
		})

//  run a transition
 err := statemachine.Emit("cook")
  if err != nil{
  	// handle error here
  	panic(err)
  }
  
//	validate
	fmt.Println(statemachine.GetCurrentState().GetIdentifier()) // outputs "cooked"
}
```