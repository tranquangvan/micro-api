ARG GO_VERSION=1.19

FROM golang:${GO_VERSION}-alpine AS builder


RUN mkdir /app
RUN ln -s /Users/tranquangvan/Desktop/micro-api app

WORKDIR /app

ADD go.mod .
ADD go.sum .

RUN go mod download
COPY . /app/

RUN go install -mod=mod github.com/githubnemo/CompileDaemon
RUN go get github.com/pilu/fresh

ENTRYPOINT CompileDaemon --build="go build main.go" --command=./main
