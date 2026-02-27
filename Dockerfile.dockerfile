FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o schedule-bot main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates chromium chromium-chromedriver

WORKDIR /root/
COPY --from=builder /app/schedule-bot .
COPY .env .

# Устанавливаем переменные для Chrome
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/lib/chromium/

CMD ["./schedule-bot"]