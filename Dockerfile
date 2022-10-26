# syntax = docker/dockerfile:1.4

# frontend
FROM node:19 as build-node

COPY . /build
WORKDIR /build
ENV NODE_ENV=production
RUN \
    --mount=type=cache,target=/build/cmd/httpserver/public/node_modules \
    make node-build

# backend
FROM golang:latest as build-go

RUN \
    --mount=type=cache,target=/var/cache/apt \
    --mount=type=cache,target=/var/lib/apt \
    apt-get update && apt-get install --assume-yes upx
COPY . /build
COPY --from=build-node /build/cmd/httpserver/public/dist/ /build/cmd/httpserver/public/dist/
WORKDIR /build
RUN \
    --mount=type=cache,target=/root/.cache \
    --mount=type=cache,target=/go \
    make go-build
RUN upx --best --lzma httpserver

# runtime
FROM alpine:latest

RUN apk add --no-cache ca-certificates
RUN [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf
COPY --from=build-go /build/httpserver /app/httpserver

EXPOSE 8080
WORKDIR /app
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
CMD ["/app/httpserver"]
