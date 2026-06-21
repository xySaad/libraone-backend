FROM golang:1.26-bookworm AS builder
RUN apt-get update && apt-get install -y --no-install-recommends \
        gcc \
        libc6-dev \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /src
ENV CGO_ENABLED=1 \
    GOOS=linux
# Cache deps separately from source for faster rebuilds.
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN go build -trimpath -ldflags="-s -w" -o /out/libraone .
RUN go build -trimpath -tags sqlite3 -ldflags="-s -w" \
    -o /out/migrate github.com/golang-migrate/migrate/v4/cmd/migrate

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /out/libraone   /app/libraone
COPY --from=builder /out/migrate    /app/migrate
COPY db/migrations                  /app/db/migrations
COPY entrypoint.sh                  /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh /app/migrate /app/libraone \
    && mkdir -p /app/db
# Holds db/database.db — mount a volume here so data survives rebuilds.
VOLUME ["/app/db"]

EXPOSE 5051

ENTRYPOINT ["/app/entrypoint.sh"]