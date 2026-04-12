FROM dev-registry.abr.ru/all_image/gi/golang-swag:1.25.3 AS builder
WORKDIR /src
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /src/out/trash_project ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/out/migrator ./cmd/migrator

FROM scratch
WORKDIR /app
# Копируем оболочку из builder
COPY --from=builder /bin/sh /bin/sh
# Копируем необходимые библиотеки (если нужны)
COPY --from=builder /lib/*/libc.so.* /lib/
COPY --from=builder /lib/*/ld-linux-* /lib/
COPY --from=builder /lib/ld-musl-* /lib/
COPY --from=builder /src/out/trash_project ./trash_project
COPY --from=builder /src/out/migrator ./migrator
COPY --from=builder /src/docs ./docs
COPY --from=builder /src/data ./data
COPY --chmod=755 docker-entrypoint.sh docker-entrypoint.sh
EXPOSE 8001
CMD ["./docker-entrypoint.sh"]
