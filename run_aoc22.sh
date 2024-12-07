#!/usr/bin/env bash

go test ./aoc22/...

for d in aoc22/day* ; do
  echo "################################################################################"
  if [[ "$d" = "aoc22/day16" ]]; then
    echo "skipping day16"
  else
    echo "RUNNING: ${d}"
    if ! go run "github.com/phicode/challenges/${d}" ; then
      echo "ABNORMAL TERMINATION"
      exit 1
    fi
  fi
done
echo "################################################################################"
