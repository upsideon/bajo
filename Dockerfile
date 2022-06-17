# syntax=docker/dockerfile:1

FROM golang:1.18-alpine AS builder

RUN mkdir /app
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod verify
COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o /bajo
CMD ["/bajo"]
