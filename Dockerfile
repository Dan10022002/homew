# syntax=docker/dockerfile:1

FROM golang:latest
WORKDIR /cmd
COPY / .

RUN go mod tidy
RUN go build -o main cmd/main.go
CMD ["./main"]
