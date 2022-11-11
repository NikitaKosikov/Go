FROM golang:alpine as builder

COPY . /test/
WORKDIR /test/

RUN go mod download
RUN go build -o ./.bin/app ./cmd/app/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /test/.bin/app .
COPY --from=builder /test/configs configs/

EXPOSE 80

CMD ["./app"] 