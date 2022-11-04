FROM golang:alpine

WORKDIR /
ADD  . .
RUN go build -o test cmd/app/main.go
CMD ["./cmd/app/main"] 