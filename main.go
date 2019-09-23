package main

import (
	"fmt"

	"github.com/deolu-asenuga/state-machine/lib"
)

type State struct {
	status string
}

func main() {
	sm := lib.NewMachine("solid")

	// define the different transitions
	sm.AddTransition(&lib.Transition{
		Name: "melt",
		From: "solid",
		To:   "liquid",
		Action: func(state interface{}) {
			fmt.Println("melted!")
		},
	})
	sm.AddTransition(&lib.Transition{
		Name: "freeze",
		From: "liquid",
		To:   "solid",
		Action: func(state interface{}) {
			fmt.Println("freezed!")
		},
	})
	sm.AddTransition(&lib.Transition{
		Name: "freeze",
		From: "liquid",
		To:   "solid",
		Action: func(state interface{}) {
			fmt.Println("frozen!")
		},
	})

	sm.AddTransition(&lib.Transition{
		Name: "condense",
		From: "gas",
		To:   "liquid",
		Action: func(state interface{}) {
			fmt.Println("condensed!")
		},
	})
	sm.AddTransition(&lib.Transition{
		Name: "evapourate",
		From: "liquid",
		To:   "gas",
		Action: func(state interface{}) {
			fmt.Println("evapourated!")
		},
	})

	// running single transitions
	sm.RunTransition("melt", func(state interface{}) {
		// the call back function
		value, ok := state.(string)
		if ok {
			fmt.Println("new state ", value)
		}
	})
	sm.RunTransition("freeze", func(state interface{}) {
		// the call back function
		value, ok := state.(string)
		if ok {
			fmt.Println("new state ", value)
		}
	})

	// running series of transitions
	sm.RunTransitionSequence([]string{"melt", "evapourate", "condense", "freeze"})
}
