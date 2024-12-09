package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseBlocks(t *testing.T) {
	d := ParseBlocks("12345")
	s := d.String()
	assert.Equal(t, "0..111....22222", s)
}

func TestParseBlocks2(t *testing.T) {
	d := ParseBlocks("2333133121414131402")
	s := d.String()
	assert.Equal(t, "00...111...2...333.44.5555.6666.777.888899", s)
}

func TestBlocksSolvePart1(t *testing.T) {
	d := ParseBlocks("2333133121414131402")
	Rearrange(d)
	s := d.String()
	assert.Equal(t, "00992111777.44.333....5555.6666.....8888..", s)
}
