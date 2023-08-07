#!/bin/sh
set -e

BUILD_DOCKER_IMG=0
DIRS="";
GO_VERSION="";

while [ "$#" -gt 0 ]; do
  case "$1" in
    -b|--build) BUILD_DOCKER_IMG=1; shift 1;;
    -gv|--go-version) GO_VERSION="$2"; shift 2;;

    -*) echo "unknown option: $1" >&2; exit 1;;
    *) DIRS="$DIRS $1"; shift 1;;
  esac
done

if [ "$DIRS" = "" ]; then
  DIRS="./dbfixtures";
fi

if [ $BUILD_DOCKER_IMG -eq 1 ]; then
  echo "Creating Dockerfile from template";
  sed "s/%%GO_VERSION%%/$GO_VERSION/g" ./docker/Dockerfile > ./docker/Dockerfile-$GO_VERSION;

  echo "Build the Docker image";
  docker build -f ./docker/Dockerfile-$GO_VERSION --pull --rm -t go-dbfixtures-$GO_VERSION:latest .;
fi

mkdir -p ./coverage/;

docker run --rm -v "${PWD}/":"/usr/src/app/" go-dbfixtures-$GO_VERSION:latest /bin/sh -c "go test -coverprofile coverage/coverage.out $DIRS && go tool cover -html coverage/coverage.out -o coverage/coverage.html && gcov2lcov -infile=coverage/coverage.out -outfile=coverage/coverage.lcov";