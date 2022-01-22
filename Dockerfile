FROM golang:1.17-alpine as builder
RUN apk update && apk add --no-cache git ca-certificates tzdata openssh make
WORKDIR /app
COPY go.mod go.sum Makefile ./
RUN go mod download
COPY . .
RUN make build
ENTRYPOINT ["/app/build/folder_manager"]

