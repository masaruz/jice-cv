#!/bin/bash
for i in $(seq 1 56)
do
  if (( $i < 10 ))
  then
    i="0$i"
  fi
  echo "Running TestLoop$i ..."
  env=dev go test -run "TestLoop$i"
done