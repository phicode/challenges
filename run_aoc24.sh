#!/usr/bin/env bash

go test aoc24/...

for d in aoc24/day* ; do
  echo "################################################################################"
  echo "RUNNING: ${d}"
  if ! go run "git.bind.ch/phil/challenges/${d}" ; then
    echo "ABNORMAL TERMINATION"
    exit 1
  fi
  echo "################################################################################"
done
