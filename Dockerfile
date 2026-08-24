FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o gin-app ./cmd/server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/gin-app .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/config.yaml .

EXPOSE 2808

CMD ["./gin-app"]