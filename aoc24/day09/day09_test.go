package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseInput(t *testing.T) {
	i := ParseInput("12345")
	s := i.String()
	assert.Equal(t, "0..111....22222", s)
}
