# gpagdispo

`gpagdispo` (_página disponible_ available website in Spanish) is a
Go system to monitor websites periodically in a scalable way.

It comprises of two Go applications: `gpagdispo-checker` and
`gpagdispo-recorder`.

### gpagdispo-checker

`gpagdispo-checker` is a Go app that reads from a JSON or
[ION](https://amzn.github.io/ion-docs/docs/spec.html) formatted files the websites to monitor
optionally matching a regular expression. It follows this format:

```json
{
  "websites": [
    {
      "url": "http://awesome.web.com",
      "method": "HEAD"
    },
    {
      "url": "https://another.awesome.web.com/placebo",
      "match_regex": "tumbles?"
    }
  ]
}
```

It periodically checks the availability of those websites,
configurable via `TICK_TIME` environment variable and send to a Kafka
topic `website.monitor` through broker configurable via
`KAFKA_ADDRS` the result of the monitor check.

### pagdispo-recorder

`pagdispo-recorder` is a Go app that reads from a `website.monitor` Kafka topic through a Kafka
broker via `KAFKA_ADDRS` the results of monitor checks of websites
and stores them in a PostgreSQL database whose DSN is configurable via
`POSTGRESQL_DSN` environment variable.

## Development

It provides a Docker compose with a Kafka + PostgreSQL ready to be
use.

```shell
make start-dev
```

Then, run `cd checker && go run ./cmd/gpagdispo-checker/main.go` for local
testing in one terminal and `cd recorder &&
go run ./cmd/gpagdipso-recorder/main.go` in other terminal.

In order to stop the docker compose, run `make stop-dev`.

## Testing

You can run all test suite running:

```shell
make start-dev
make test
make integration-test
```

### CI

We are using [Github Actions](.github/workflows/test.yml) for
running the tests and linting in parallel.

## Linting

[golangci-lint](https://golangci-lint.run) is used for linting the Go
code and [shellcheck](https://www.shellcheck.net/) for Shell script linting.

You can run
```
make lint
```

To see their results.

## Settings

The configuration settings can be modified using environment variables.

For example:

```
KAFKA_ADDRS=1.1.1.1:9092 KAFKA_CERT_FILE=service.cert KAFKA_KEY_FILE=service.key KAFKA_CA_FILE=ca.pem POSTGRESQL_DSN=postgres://admin:pass@host:20127/defaultdb?sslmode=require go run ./cmd/gpagdispo-recorder/main.go
```

Available settings for
[gpagdispo-checker](checker/cmd/gpagdispo-checker/main.go) and [gpagdispo-recorder](recorder/cmd/gpagdispo-recorder/main.go).
