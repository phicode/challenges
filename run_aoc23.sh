#!/usr/bin/env bash

go test aoc23/...

for d in aoc23/day* ; do
  echo "################################################################################"
  echo "RUNNING: ${d}"
  go run git.bind.ch/phil/challenges/${d}
  if [[ $? -ne 0 ]]; then
    echo "ABNORMAL TERMINATION"
    exit 1
  fi
  echo "################################################################################"
done

