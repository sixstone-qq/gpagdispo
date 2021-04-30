# First stage
FROM golang:1.16-alpine as builder

RUN apk --no-cache add git ca-certificates musl-dev

ADD . /go/src/github.com/sixstone-qq/gpagdispo/checker
WORKDIR /go/src/github.com/sixstone-qq/gpagdispo/checker

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -a -installsuffix cgo ./cmd/gpagdispo-checker/

# Second stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY websites.ion /app/websites.ion
ENV CONFIG_PATH=websites.ion

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /go/src/github.com/sixstone-qq/gpagdispo/checker/gpagdispo-checker .

CMD ["/app/gpagdispo-checker"]
