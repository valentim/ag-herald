FROM golang:1.14-alpine

ENV SRC_PATH $GOPATH/src/github.com/valentim/ag-herald
ENV GO111MODULE=on

WORKDIR $SRC_PATH

RUN apk add --update git gcc libc-dev

COPY . $SRC_PATH

RUN go get -d -v

RUN go build main.go

EXPOSE 4000