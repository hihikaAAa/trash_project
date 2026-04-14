FROM golang:1.25-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /src/out/trash_project ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/out/migrator ./cmd/migrator

FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache ca-certificates && update-ca-certificates

COPY --from=builder /src/out/trash_project ./trash_project
COPY --from=builder /src/out/migrator ./migrator
COPY --from=builder /src/data ./data
COPY --chmod=755 docker-entrypoint.sh ./docker-entrypoint.sh

EXPOSE 8001
CMD ["./docker-entrypoint.sh"]
