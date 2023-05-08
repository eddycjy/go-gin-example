FROM golang:1.20

RUN go install github.com/cosmtrek/air@latest

WORKDIR /app/
COPY . .
RUN go mod vendor

CMD ["air"]