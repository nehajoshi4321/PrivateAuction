#!/bin/bash

declare -A input=([0,0]=a [0,1]=b [1,0]=c [1,1]=d)

   #(5 1024 1000000000 10 10000 100000 76)
   # (5 1024 1000000000 15 10000 100000 76)
   #(5 1024 1000000000 20 10000 100000 76)
   # (5 2048 1000000000 5 10000 100000 76)

for i in "${input[@]}"
do
  # Run your Go program with the input value as a single argument
  echo ${{$i}[0]};
  go run test1.go ${{$i}[0]}
done
