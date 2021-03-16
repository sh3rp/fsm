package fsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThreeStates(t *testing.T) {
	s1 := State(1)
	s2 := State(2)
	s3 := State(3)

	fsm := NewFSM()

	fsm.RegisterState(s1, nil, nil)
	fsm.RegisterState(s2, nil, nil)
	fsm.RegisterState(s3, nil, nil)

	fsm.RegisterTransition(s1, s2)
	fsm.RegisterTransition(s2, s3)

	err := fsm.Initialize(s1)
	assert.Nil(t, err)

	err = fsm.Transition(s2, nil)
	assert.Nil(t, err)

	err = fsm.Transition(s3, nil)
	assert.Nil(t, err)

	assert.Equal(t, fsm.Current(), s3)
}

func TestStatesWithTransitionFuncs(t *testing.T) {
	s1 := State(1)
	s2 := State(2)

	fsm := NewFSM()

	var s1Metadata, s2Metadata map[string]string
	s1Left := false
	s2Entered := false

	fsm.RegisterState(s1, nil, func(s State, m map[string]string) {
		s1Left = true
		s1Metadata = m
	})
	fsm.RegisterState(s2, func(s State, m map[string]string) {
		s2Entered = true
		s2Metadata = m
	}, nil)
	fsm.RegisterTransition(s1, s2)

	fsm.Initialize(s1)
	fsm.Transition(s2, map[string]string{"key1": "value1"})

	assert.True(t, s1Left)
	assert.True(t, s2Entered)
	assert.Equal(t, s1Metadata["key1"], "value1")
	assert.Equal(t, s2Metadata["key1"], "value1")
}
