FROM golang:1.16-alpine

RUN apk --no-cache add curl

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s

ENTRYPOINT air
