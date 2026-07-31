FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o go-gin-example .

FROM alpine:3.20
RUN adduser -D -u 1000 appuser
WORKDIR /app

COPY --from=builder /app/go-gin-example .
COPY --from=builder /app/conf ./conf

RUN mkdir -p runtime/fonts runtime/qrcode upload/images export qrcode logs && \
    chown -R appuser:appuser /app




USER appuser
EXPOSE 8000
ENTRYPOINT ["./go-gin-example"]