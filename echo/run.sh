#!/bin/bash

set -x

mkdir -p logs

go build barrage.go

ps -ef | grep barrage  | grep -v grep | awk '{system("kill -9 "$2)}'
nohup ./barrage > logs/barrage.log 2>&1 &

tail -f logs/barrage.log
