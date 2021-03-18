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

	fsm.RegisterState(s1)
	fsm.RegisterState(s2)
	fsm.RegisterState(s3)

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

	fsm.RegisterState(s1)
	fsm.Leave(s1, func(s State, m map[string]string) {
		s1Left = true
		s1Metadata = m
	})
	fsm.RegisterState(s2)
	fsm.Enter(s2, func(s State, m map[string]string) {
		s2Entered = true
		s2Metadata = m
	})
	fsm.RegisterTransition(s1, s2)

	fsm.Initialize(s1)
	fsm.Transition(s2, map[string]string{"key1": "value1"})

	assert.True(t, s1Left)
	assert.True(t, s2Entered)
	assert.Equal(t, s1Metadata["key1"], "value1")
	assert.Equal(t, s2Metadata["key1"], "value1")
}

func TestThreeStatesBroken(t *testing.T) {
	s1 := State(1)
	s2 := State(2)
	s3 := State(3)
	s4 := State(4)

	fsm := NewFSM()

	fsm.RegisterState(s1)
	fsm.RegisterState(s2)
	fsm.RegisterState(s3)

	fsm.RegisterTransition(s1, s2)
	fsm.RegisterTransition(s2, s3)

	err := fsm.Initialize(s1)
	assert.Nil(t, err)

	err = fsm.Transition(s3, nil)
	assert.NotNil(t, err)
	assert.Equal(t, "cannot transition to this state", err.Error())

	fsm = NewFSM()

	fsm.RegisterState(s1)
	fsm.RegisterState(s2)
	fsm.RegisterState(s3)

	fsm.RegisterTransition(s1, s2)
	fsm.RegisterTransition(s2, s3)

	err = fsm.Initialize(s4)
	assert.NotNil(t, err)

	fsm.RegisterState(s4)
	err = fsm.Initialize(s4)
	assert.Nil(t, err)

	err = fsm.Transition(s1, nil)
	assert.NotNil(t, err)
	assert.Equal(t, "cannot transition from this state", err.Error())
}
