#!/bin/sh

if [ ! -f hel ]; then
    go build .
fi

./hel \
    --input=http://m.dk/structure/transit.ashx/status \
    --endpoint=/ \
    --port=9000 \
    --service=badrequest \
    --delay=0

