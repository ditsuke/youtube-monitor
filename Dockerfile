FROM golang:1.19-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
COPY --from=builder /app/main /app/main
CMD ["/app/main"]
