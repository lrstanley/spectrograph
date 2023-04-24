# syntax = docker/dockerfile:1.4

# frontend
FROM node:20 as build-node

COPY . /build
WORKDIR /build
ENV NODE_ENV=production
RUN \
    --mount=type=cache,target=/build/cmd/httpserver/public/node_modules \
    make node-build

# backend
FROM golang:latest as build-go

COPY . /build
COPY --from=build-node /build/cmd/httpserver/public/dist/ /build/cmd/httpserver/public/dist/
WORKDIR /build
RUN \
    --mount=type=cache,target=/root/.cache \
    --mount=type=cache,target=/go \
    make go-build

# runtime
FROM alpine:3.17

RUN apk add --no-cache ca-certificates
RUN if [ ! -e /etc/nsswitch.conf ];then echo 'hosts: files dns' > /etc/nsswitch.conf;fi
COPY --from=build-go /build/httpserver /app/httpserver

EXPOSE 8080
WORKDIR /app
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
CMD ["/app/httpserver"]
