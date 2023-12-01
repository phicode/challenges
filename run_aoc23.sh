#!/usr/bin/env bash

go test aoc23/...

for d in aoc23/day* ; do
  echo "################################################################################"
  echo "RUNNING: ${d}"
  go run git.bind.ch/phil/challenges/${d}
  echo "################################################################################"
done

