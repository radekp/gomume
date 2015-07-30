#!/bin/sh
stty -icanon min 1 time 0
go run main.go
stty cooked echo ctlecho
echo done
