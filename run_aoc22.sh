#!/usr/bin/env bash

go test aoc22/...

for d in aoc22/day* ; do
  if [[ "$d" = "aoc22/day16" ]]; then
    echo "skipping day16"
  else
    echo "################################################################################"
    echo "RUNNING: ${d}"
    go run git.bind.ch/phil/challenges/${d}
    echo "################################################################################"
  fi
done

