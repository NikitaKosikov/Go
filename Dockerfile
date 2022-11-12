FROM golang:alpine as builder

COPY . /test/
WORKDIR /test/

RUN go mod download
RUN go build -o ./.bin/app ./cmd/app/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/zhashkevych/telegram-pocket-bot/.bin/bot .
COPY --from=0 /github.com/zhashkevych/telegram-pocket-bot/configs configs/

EXPOSE 4000

CMD ["./app"]