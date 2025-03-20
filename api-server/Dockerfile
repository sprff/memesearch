FROM golang:1.24.0-alpine3.20 AS builder
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY cmd/apiserver ./cmd/app
COPY internal ./internal

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o /app/bin/app ./cmd/app


FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/bin/app /app/bin/app
COPY config ./config
EXPOSE 1781
WORKDIR /app


ENTRYPOINT ["/app/bin/app"]
