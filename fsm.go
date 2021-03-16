package fsm

import "fmt"

type State int

type Transition struct {
	From State
	To   State
}

type stateWrapper struct {
	state State
	enter func(State, map[string]string)
	leave func(State, map[string]string)
}

type FSM interface {
	RegisterState(State, func(State, map[string]string), func(State, map[string]string))
	RegisterTransition(State, State)
	Current() State
	Initialize(State) error
	Transition(State, map[string]string) error
}

type fsm struct {
	validStates      map[State]stateWrapper
	validTransitions map[State][]Transition
	transitions      []Transition
	currentState     State
}

func NewFSM() FSM {
	return &fsm{validStates: make(map[State]stateWrapper), validTransitions: make(map[State][]Transition), transitions: []Transition{}, currentState: State(0)}
}

func (fsm *fsm) RegisterState(state State, enterFunc func(State, map[string]string), exitFunc func(State, map[string]string)) {
	fsm.validStates[state] = stateWrapper{state, enterFunc, exitFunc}
}

func (fsm *fsm) RegisterTransition(from State, to State) {
	if _, hasTransitions := fsm.validTransitions[from]; !hasTransitions {
		fsm.validTransitions[from] = []Transition{}
	}
	fsm.validTransitions[from] = append(fsm.validTransitions[from], Transition{from, to})
}

func (fsm *fsm) Current() State {
	return fsm.currentState
}

func (fsm *fsm) Initialize(state State) error {
	if _, hasState := fsm.validStates[state]; !hasState {
		return fmt.Errorf("state does not exist")
	}

	fsm.currentState = state

	return nil
}

func (fsm *fsm) Transition(to State, metadata map[string]string) error {
	if _, isValid := fsm.validTransitions[fsm.currentState]; !isValid {
		return fmt.Errorf("cannot transition from this state")
	}

	vt := fsm.validTransitions[fsm.currentState]

	found := false
	for _, trans := range vt {
		if trans.To == to {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("cannot transition to this state")
	}

	fsm.transitions = append(fsm.transitions, Transition{fsm.currentState, to})
	if fsm.validStates[fsm.currentState].leave != nil {
		fsm.validStates[fsm.currentState].leave(to, metadata)
	}
	from := fsm.currentState
	fsm.currentState = to
	if fsm.validStates[fsm.currentState].enter != nil {
		fsm.validStates[fsm.currentState].enter(from, metadata)
	}

	return nil
}
