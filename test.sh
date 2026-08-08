#!/bin/bash

go build -o app . || exit 1

cleanup() {
    echo "Stopping applications..."

    echo $PID_T $PID_R
    kill -TERM "$PID_T" "$PID_R"

    wait "$PID_T" "$PID_R" 2>/dev/null

    echo "Stopped."
    exit 0
}

trap cleanup SIGINT SIGTERM SIGKILL

./app -n a -m t &
PID_T=$!

./app -n b -m r &
PID_R=$!

wait