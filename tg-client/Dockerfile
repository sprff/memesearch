FROM golang:1.24.0-alpine3.20 AS builder
RUN apk add build-base
WORKDIR /api-client
COPY api-client/go.mod api-client/go.sum /api-client/
RUN go mod download
WORKDIR /app
COPY tg-client/go.mod tg-client/go.sum /app/
RUN go mod download


COPY api-client/internal /api-client/internal
COPY api-client/pkg /api-client/pkg

COPY tg-client/cmd/main /app/cmd/app
COPY tg-client/internal /app/internal

ENV GOCACHE=/root/.cache/go-build
ENV CGO_ENABLED=1
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o /app/bin/app ./cmd/app


FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/bin/app /app/bin/app

EXPOSE 1781
WORKDIR /app


ENTRYPOINT ["/app/bin/app"]
