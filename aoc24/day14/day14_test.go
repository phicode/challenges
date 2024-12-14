package main

import (
	"testing"

	"github.com/phicode/challenges/lib/rowcol"
	"github.com/stretchr/testify/assert"
)

func TestRobot_Move(t *testing.T) {
	r := Robot{
		P: rowcol.Pos{1, 1},
		V: rowcol.Pos{-1, -1},
	}
	move1 := r.Move(5, 5)
	assert.Equal(t, 0, move1.P.Row)
	assert.Equal(t, 0, move1.P.Col)

	move2 := move1.Move(5, 5)
	assert.Equal(t, 4, move2.P.Row)
	assert.Equal(t, 4, move2.P.Col)
}
