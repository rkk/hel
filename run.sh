#!/bin/sh

binary="hel"
if [ -f "${binary}" ]; then
    rm "${binary}"
fi

if [ ! -f "${binary}" ]; then
    output=$(go build .)
    build_res="${?}"
    if [ "${build_res}" -ne 0 ]; then
        echo "ERROR: Cannot build"
        echo "${output}"
        echo ""
        exit "${build_res}"
    fi
fi

./${binary} \
    --input=http://m.dk/structure/transit.ashx/status \
    --endpoint=/ \
    --port=9000 \
    --service=badrequest \
    --delay=0

