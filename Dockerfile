FROM golang:1.19-alpine3.16 as builder
WORKDIR /app

COPY go.* .
RUN go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build go build -o server .
RUN --mount=type=cache,target=/root/.cache/go-build go build -o migrator cmd/generate/main.go

FROM alpine:3.16
COPY --from=builder /app/server /app/migrator /app/docker/start.sh /app/.env /app/

# dotenv needs this to load the env file (not that we need it here)
WORKDIR /app

CMD ["/app/start.sh"]
