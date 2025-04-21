FROM golang:1.24.1-bookworm AS builder

RUN apt-get update && apt-get install -y gcc libsqlite3-dev sqlite3

WORKDIR /app

COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o forum ./main
RUN ls -l /app

FROM golang:1.24.1-bookworm
WORKDIR /app

RUN apt-get update && apt-get install -y sqlite3 libsqlite3-0

COPY --from=builder /app/forum .
COPY --from=builder /app/main ./main
COPY --from=builder /app/static ./static

RUN chmod +x /app/forum

EXPOSE 8080

CMD [ "/app/forum" ]