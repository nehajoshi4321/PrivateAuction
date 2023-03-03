#!/bin/bash

values=(
    ("5" "512" "1000000000" "5" "10000" "100000" "76")
    ("5" "1024" "1000000000" "5" "10000" "100000" "76")
    ("5" "1024" "1000000000" "10" "10000" "100000" "76")
    ("5" "1024" "1000000000" "15" "10000" "100000" "76")
    ("5" "1024" "1000000000" "20" "10000" "100000" "76")
    ("5" "2048" "1000000000000" "5" "10000" "100000" "76"))
    
for i in "${values[@]}"
do
    go run test1.go "${i[@]}"
done
