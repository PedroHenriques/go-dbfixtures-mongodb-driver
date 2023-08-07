#!/bin/sh
set -e

PULL_FLAG="";

while [ "$#" -gt 0 ]; do
  case "$1" in
    --pull) PULL_FLAG="--pull always"; shift 1;;

    -*) echo "unknown option: $1" >&2; exit 1;;
    *) DIRS="$DIRS $1"; shift 1;;
  esac
done

docker run --rm -v "${PWD}/":"/usr/src/app/" $PULL_FLAG -w "/usr/src/app" golangci/golangci-lint:latest /bin/sh -c "go mod tidy && golangci-lint run -v";