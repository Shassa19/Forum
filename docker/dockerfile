FROM golang:1.20-alpine AS builder 

WORKDIR /app

COPY go.mod go.sum ./
RUN go env -w GO111MODULE=on
RUN go mod download
COPY . .
RUN go build -o forum .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/forum .

CMD ["./forum"]