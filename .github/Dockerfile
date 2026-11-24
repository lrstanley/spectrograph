# syntax=docker/dockerfile:1

# frontend (Nuxt SPA, see web/).
FROM oven/bun:1 AS build-web
WORKDIR /build/web
COPY web/package.json web/bun.lock ./
RUN --mount=type=cache,target=/root/.bun/install/cache \
	bun install --frozen-lockfile
COPY web/ ./
# The frontend API client is generated from the backend OpenAPI spec.
COPY internal/database/ent/rest/openapi.json /build/internal/database/ent/rest/openapi.json
RUN bun --bun run generate

# backend.
FROM golang:1.27 AS build-go
WORKDIR /build
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
	go mod download
COPY . .
COPY --from=build-web /build/web/.output/public/ /build/cmd/httpserver/public/dist/
RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 go build \
		-ldflags '-d -s -w -extldflags=-static' \
		-tags=netgo,osusergo,static_build \
		-installsuffix netgo \
		-trimpath \
		-o /build/spectrograph \
		./cmd/httpserver

# runtime.
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build-go /build/spectrograph /spectrograph
ENV HTTP_BIND_ADDR=":8080"
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/spectrograph"]
