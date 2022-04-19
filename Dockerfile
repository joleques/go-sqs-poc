FROM golang:1.17.0-alpine3.14 as builder

WORKDIR /go/src/github.com/joleques/go-sqs-poc

COPY . .

RUN go mod download

RUN go build -ldflags "-s -w" src/main.go

FROM alpine:3.14

WORKDIR /app

COPY --from=builder /go/src/github.com/joleques/go-sqs-poc/main .

CMD [ "./main" ]
