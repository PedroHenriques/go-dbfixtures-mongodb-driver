FROM golang:%%GO_VERSION%%-buster

WORKDIR /usr/src/app

# Copy the repo into the container
ADD . .

# Install GO dependencies
RUN go mod tidy

# Install dependency needed for converting coverage output into lcov format
RUN go install github.com/jandelgado/gcov2lcov@latest

# Install dependency needed for running tests in watch mode
RUN go install github.com/mitranim/gow@latest || true