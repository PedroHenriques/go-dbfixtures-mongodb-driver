#!/bin/sh
set -e

BUILD_DOCKER_IMG=0
WATCH=0;
DIRS="";
GO_VERSION="";

while [ "$#" -gt 0 ]; do
  case "$1" in
    -b|--build) BUILD_DOCKER_IMG=1; shift 1;;
    -w|--watch) WATCH=1; shift 1;;
    -gv|--go-version) GO_VERSION="$2"; shift 2;;

    -*) echo "unknown option: $1" >&2; exit 1;;
    *) DIRS="$DIRS $1"; shift 1;;
  esac
done

if [ "$DIRS" = "" ]; then
  DIRS="./...";
fi

if [ $BUILD_DOCKER_IMG -eq 1 ]; then
  echo "Creating Dockerfile from template";
  sed "s/%%GO_VERSION%%/$GO_VERSION/g" ./docker/Dockerfile > ./docker/Dockerfile-$GO_VERSION;

  echo "Build the Docker image";
  docker build -f ./docker/Dockerfile-$GO_VERSION --pull --rm -t go-dbfixtures-$GO_VERSION:latest .;
fi

CMD="go test -v -cover $DIRS";
DOCKER_FLAGS="";
if [ $WATCH -eq 1 ]; then
  CMD="gow -c test -v -cover $DIRS";
  DOCKER_FLAGS="-it"; # Can not be always added since these docker flags are not supported in Github actions
fi

docker network create tests || true;
docker run --rm --network tests --name testmongo -p 127.0.0.1:27017:27017/tcp -d mongo:4 || true;
docker run --rm --network tests $DOCKER_FLAGS -v "${PWD}/":"/usr/src/app/" go-dbfixtures-$GO_VERSION:latest /bin/sh -c "$CMD";