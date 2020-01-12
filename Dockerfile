#build stage
FROM golang:alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN apk add --no-cache git &&\
    go mod download &&\
    go build -o telegramNotifier .

#final stage
FROM alpine:latest
LABEL Name=telegramnotifier Version=0.0.1
EXPOSE 8080
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/app/telegramNotifier /app/telegramNotifier
ENTRYPOINT ./telegramNotifier
