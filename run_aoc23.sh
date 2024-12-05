#!/usr/bin/env bash

go test aoc23/...

for d in aoc23/day* ; do
  echo "################################################################################"
  echo "RUNNING: ${d}"
  if ! go run "git.bind.ch/phil/challenges/${d}" ; then
    echo "ABNORMAL TERMINATION"
    exit 1
  fi
done
echo "################################################################################"
