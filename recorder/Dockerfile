# First stage
FROM golang:1.16-alpine as builder

RUN apk --no-cache add git ca-certificates musl-dev

ADD . /go/src/github.com/sixstone-qq/gpagdispo/recorder
WORKDIR /go/src/github.com/sixstone-qq/gpagdispo/recorder

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -a -installsuffix cgo ./cmd/gpagdispo-recorder/

# Second stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY db /app/db

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /go/src/github.com/sixstone-qq/gpagdispo/recorder/gpagdispo-recorder .

CMD ["/app/gpagdispo-recorder"]
