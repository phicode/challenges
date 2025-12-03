#!/usr/bin/env bash

go test ./aoc25/...

for d in aoc25/day* ; do
  echo "################################################################################"
  echo "RUNNING: ${d}"
  if ! go run "github.com/phicode/challenges/${d}" ; then
    echo "ABNORMAL TERMINATION"
    exit 1
  fi
done
echo "################################################################################"
