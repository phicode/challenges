package rowcol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPos_Right(t *testing.T) {
	assert.Equal(t, Right, Up.Right())
	assert.Equal(t, Down, Right.Right())
	assert.Equal(t, Left, Down.Right())
	assert.Equal(t, Up, Left.Right())
}

func TestPos_Left(t *testing.T) {
	assert.Equal(t, Left, Up.Left())
	assert.Equal(t, Down, Left.Left())
	assert.Equal(t, Right, Down.Left())
	assert.Equal(t, Up, Right.Left())
}
