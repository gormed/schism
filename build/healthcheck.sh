#!/bin/sh
set -e

curl --noproxy 127.0.0.1 --max-time 2 --fail "http://127.0.0.1:8080" || ( echo "healthcheck failed" && exit 1 )
