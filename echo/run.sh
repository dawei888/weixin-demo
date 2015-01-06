#!/bin/bash

set -x

mkdir -p logs

go build echo.go

ps -ef | grep echo  | grep -v grep | awk '{system("kill -9 "$2)}'
nohup ./echo > logs/echo.log 2>&1 &

tail -f logs/echo.log
