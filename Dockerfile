# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-buster AS build

WORKDIR /

COPY / ./

RUN go mod download &&\
    go build -o /sticker

## Deploy
FROM alpine:3.17.0

WORKDIR /

COPY --from=build /sticker /sticker
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

EXPOSE 8080

ENTRYPOINT ["/sticker"]