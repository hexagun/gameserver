package gameserver

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type State struct {
	// Define your state properties
	Value int
}

type IncrementPayload struct {
}

type IncrementAction struct {
	Type    string
	Payload IncrementPayload
}

func (r IncrementAction) GetType() string {
	return r.Type
}

func (r IncrementAction) GetPayload() interface{} {
	return r.Payload
}

type DecrementPayload struct {
}

type DecrementAction struct {
	Type    string
	Payload DecrementPayload
}

func (r DecrementAction) GetType() string {
	return r.Type
}

func (r DecrementAction) GetPayload() interface{} {
	return r.Payload
}

type ActionHandlerFunc func(state *State, action Action)

func LoggingMiddleware(next ActionHandlerFunc) ActionHandlerFunc {
	return func(state *State, action Action) {
		fmt.Printf("Action Type: %s, Payload: %v", action.GetType(), action.GetPayload())
		next(state, action) // pass to the next middleware or handler
	}
}

func ValidationMiddleware(next ActionHandlerFunc) ActionHandlerFunc {
	return func(state *State, action Action) {
		if action.GetType() == "INCREMENT" && state.Value >= 10 {
			fmt.Println("State value cannot exceed 10")
			return
		}
		next(state, action) // pass to the next middleware or handler
	}
}

func ChainMiddleware(handler ActionHandlerFunc, middlewares ...func(ActionHandlerFunc) ActionHandlerFunc) ActionHandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func FinalHandler(state *State, action Action) {
	switch action.GetType() {
	case "INCREMENT":
		state.Value++
	case "DECREMENT":
		state.Value--
	default:
		fmt.Println("Unknown action type")
	}
}

func TestMiddlewareChain(t *testing.T) {
	state := &State{Value: 0}

	handler := ChainMiddleware(FinalHandler, LoggingMiddleware, ValidationMiddleware)

	actions := []Action{
		IncrementAction{Type: "INCREMENT"},
		IncrementAction{Type: "INCREMENT"},
		DecrementAction{Type: "DECREMENT"},
		IncrementAction{Type: "INCREMENT"},
		IncrementAction{Type: "INCREMENT"},
	}

	// Dispatch actions through the middleware chain
	for _, action := range actions {
		handler(state, action)
		fmt.Printf("Current State Value: %d\n", state.Value)
	}
	assert.Equal(t, 3, state.Value, "wrong value")
}
