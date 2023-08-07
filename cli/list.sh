#!/bin/sh
set -e

while [ "$#" -gt 0 ]; do
  case "$1" in
    --version) VERSION="$2"; shift 2;;

    -*) echo "unknown option: $1" >&2; exit 1;;
    *) shift 1;;
  esac
done

docker run --rm -v "${PWD}/":"/usr/src/app/" golang:1.19-buster /bin/sh -c "GOPROXY=proxy.golang.org go list -m github.com/PedroHenriques/go-dbfixtures@${VERSION}";