FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY save-csv/* ./
RUN GOOS=linux GOARCH=arm64 go build -o bin .

FROM debian:bookworm-slim AS runtime

WORKDIR /app

COPY --from=builder /app/bin .
COPY data.csv .
COPY start-script.sh ./
RUN chmod +x start-script.sh

ENTRYPOINT ["/app/start-script.sh"]