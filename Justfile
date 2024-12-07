#!/usr/bin/env -S just -f

# just - 🤖 Just a command runner - https://github.com/casey/just
# https://just.systems/man/en/introduction.html

packages     := `go list github.com/phicode/challenges/...`
package_dirs := `go list -f '{{.Dir}}' github.com/phicode/challenges/...`

@default: help

# this help
@help:
  just --list

all: test

test:
  go test github.com/phicode/challenges/...
