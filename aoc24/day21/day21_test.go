package main

import (
	"testing"

	"github.com/phicode/challenges/lib/rowcol"
	"github.com/stretchr/testify/assert"
)

func TestNumPad_Press(t *testing.T) {
	numPad := NumPad{
		V: 'A',
		P: NumPadPosition('A'),
	}
	assert.Equal(t, rowcol.PosXY(2, 3), numPad.P)
	// 789
	// 456
	// 123
	//  0A
	//
	// 029A
	moves0 := numPad.Press('0')
	assert.Equal(t, rowcol.PosXY(1, 3), numPad.P)
	assert.Equal(t, []rune("<A"), moves0)

	moves2 := numPad.Press('2')
	assert.Equal(t, rowcol.PosXY(1, 2), numPad.P)
	assert.Equal(t, []rune("^A"), moves2)

	moves9 := numPad.Press('9')
	assert.Equal(t, rowcol.PosXY(2, 0), numPad.P)
	assert.Equal(t, []rune("^^>A"), moves9)

	movesA := numPad.Press('A')
	assert.Equal(t, rowcol.PosXY(2, 3), numPad.P)
	assert.Equal(t, []rune("vvvA"), movesA)
}

func TestDirPad_Press(t *testing.T) {
	dp := DirPad{
		V: 'A',
		P: DirPadPosition('A'),
	}
	assert.Equal(t, rowcol.PosXY(2, 0), dp.P)
	//  ^A
	// <v>
	//
	// <A^A
	moves0 := dp.Press('<')
	assert.Equal(t, rowcol.PosXY(0, 1), dp.P)
	assert.Equal(t, []rune("v<<A"), moves0)

	moves2 := dp.Press('A')
	assert.Equal(t, rowcol.PosXY(2, 0), dp.P)
	assert.Equal(t, []rune(">>^A"), moves2)

	moves9 := dp.Press('^')
	assert.Equal(t, rowcol.PosXY(1, 0), dp.P)
	assert.Equal(t, []rune("<A"), moves9)

	movesA := dp.Press('A')
	assert.Equal(t, rowcol.PosXY(2, 0), dp.P)
	assert.Equal(t, []rune(">A"), movesA)
}

func TestDirPad_MoveLeftTwice(t *testing.T) {
	dp1 := DirPad{
		V: 'A',
		P: DirPadPosition('A'),
	}
	dp2 := DirPad{
		V: 'A',
		P: DirPadPosition('A'),
	}
	out1 := dp1.Press('<')
	out2 := dp2.PressMultiple(out1)
	assert.Equal(t, 10, len(out2))
}

func TestDirPad_MoveDownTwice(t *testing.T) {
	dp1 := DirPad{
		V: 'A',
		P: DirPadPosition('A'),
	}
	dp2 := DirPad{
		V: 'A',
		P: DirPadPosition('A'),
	}
	out1 := dp1.Press('v')
	out2 := dp2.PressMultiple(out1)
	assert.Equal(t, 9, len(out2))
}

func TestX(t *testing.T) {
	var cs Combinations
	cs.AddX(rowcol.Pos{}, invalidNumPadPos, rowcol.Right, rowcol.Down, 2, 3)
	assert.Equal(t, 9, len(cs.Cs))
}
