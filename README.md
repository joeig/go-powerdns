# PowerDNS 4.x API bindings for Golang

This community project provides bindings for PowerDNS Authoritative Server.
It's not associated with the official PowerDNS product itself.

[![Test coverage](https://img.shields.io/badge/coverage-100%25-success)](https://github.com/joeig/go-powerdns/tree/master/.github/testcoverage.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/joeig/go-powerdns/v3)](https://goreportcard.com/report/github.com/joeig/go-powerdns/v3)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/joeig/go-powerdns/v3)](https://pkg.go.dev/github.com/joeig/go-powerdns/v3)

## Features

Conveniently manage

* [zones](https://github.com/joeig/go-powerdns?tab=readme-ov-file#getaddchangedelete-zones)
* [resource records](https://github.com/joeig/go-powerdns?tab=readme-ov-file#addchangedelete-resource-records)
* [cryptokeys](https://github.com/joeig/go-powerdns?tab=readme-ov-file#handle-dnssec-cryptographic-material) (DNSSEC)
* [TSIG keys](https://github.com/joeig/go-powerdns?tab=readme-ov-file#createchangedelete-tsig-keys)
* [servers](https://pkg.go.dev/github.com/joeig/go-powerdns/v3#ServersService)
* [statistics](https://github.com/joeig/go-powerdns?tab=readme-ov-file#request-server-information-and-statistics)
* [configuration](https://pkg.go.dev/github.com/joeig/go-powerdns/v3#ConfigService)

It works entirely with the Go standard library and can easily be customized.[^1]

[^1]: There is a dependency for `github.com/jarcoal/httpmock`, which is used by the test suite.

For more features, consult our [documentation](https://pkg.go.dev/github.com/joeig/go-powerdns/v3).

## Setup

```shell
go get -u github.com/joeig/go-powerdns/v3
```

```go
import "github.com/joeig/go-powerdns/v3"
```

## Usage

### Initialize the handle

```go
import (
	"github.com/joeig/go-powerdns/v3"
	"context"
)

// Let's say
// * PowerDNS Authoritative Server is listening on `http://localhost:80`,
// * the virtual host is `localhost` and
// * the API key is `apipw`.
pdns := powerdns.New("http://localhost:80", "localhost", powerdns.WithAPIKey("apipw"))

// All API interactions support a Go context, which allow you to pass cancellation signals and deadlines.
// If you don't need a context, `context.Background()` would be the right choice for the following examples.
// If you want to learn more about how context helps you to build reliable APIs, see: https://go.dev/blog/context
ctx := context.Background()
```

#### Migrate `NewClient` to `New`

If you have used `NewClient` before and want to migrate to `New`, please see the [release notes for v3.13.0](https://github.com/joeig/go-powerdns/releases/tag/v3.13.0).

### Get/add/change/delete zones

```go
zones, err := pdns.Zones.List(ctx)
zone, err := pdns.Zones.Get(ctx, "example.com")
export, err := pdns.Zones.Export(ctx, "example.com")
zone, err := pdns.Zones.AddNative(ctx, "example.com", true, "", false, "foo", "foo", true, []string{"ns.foo.tld."})
err := pdns.Zones.Change(ctx, "example.com", &zone)
err := pdns.Zones.Delete(ctx, "example.com")
```

### Add/change/delete resource records

```go
err := pdns.Records.Add(ctx, "example.com", "www.example.com", powerdns.RRTypeAAAA, 60, []string{"::1"})
err := pdns.Records.Change(ctx, "example.com", "www.example.com", powerdns.RRTypeAAAA, 3600, []string{"::1"})
err := pdns.Records.Delete(ctx, "example.com", "www.example.com", powerdns.RRTypeA)
records, err := pdns.Records.Get(ctx, "example.com", "www.example.com", powerdns.RRTypePtr(powerdns.RRTypeA))
```

### Request server information and statistics

```go
statistics, err := pdns.Statistics.List(ctx)
servers, err := pdns.Servers.List(ctx)
server, err := pdns.Servers.Get(ctx, "localhost")
```

### Handle DNSSEC cryptographic material

```go
cryptokeys, err := pdns.Cryptokeys.List(ctx)
cryptokey, err := pdns.Cryptokeys.Get(ctx, "example.com", "1337")
err := pdns.Cryptokeys.Delete(ctx, "example.com", "1337")
```

### Create/change/delete TSIG keys

```go
tsigkey, err := pdns.TSIGKeys.Create(ctx, "examplekey", "hmac-sha256", "")
tsigkey, err := pdns.TSIGKeys.Change(ctx, "examplekey.", powerdns.TSIGKey{Key: powerdns.String("newkey")})
tsigkeys, err := pdns.TSIGKeys.List(ctx)
tsigkey, err := pdns.TSIGKeys.Get(ctx, "examplekey.")
err := pdns.TSIGKeys.Delete(ctx, "examplekey.")
```

### More examples

There are several examples on [pkg.go.dev](https://pkg.go.dev/github.com/joeig/go-powerdns/v3#pkg-examples).

### Documentation

See [pkg.go.dev](https://pkg.go.dev/github.com/joeig/go-powerdns/v3) for a full reference.

## Setup

### Requirements

#### Tested PowerDNS versions

Supported versions of PowerDNS Authoritative Server ("API v1"):

* 4.7
* 4.8
* 4.9

Version 4.1, 4.2 and 4.3 are probably working fine, but are officially [end-of-life](https://repo.powerdns.com/).
Be aware that there are breaking changes in "API v1" between PowerDNS 3.x, 4.0 and 4.1.

#### Tested Go versions

In accordance with [Go's version support policy](https://golang.org/doc/devel/release.html#policy), this module is being tested with the following Go releases:

* 1.22
* 1.23

## Contribution

This API client has not been completed yet, so feel free to contribute.
The [OpenAPI specification](https://github.com/PowerDNS/pdns/blob/master/docs/http-api/swagger/authoritative-api-swagger.yaml) is a good reference.

You can use Docker Compose to launch a PowerDNS authoritative server including a generic SQLite3 backend, DNSSEC support and some optional fixtures:

```bash
docker-compose -f docker-compose-v4.9.yml up
docker-compose -f docker-compose-v4.9.yml exec powerdns sh init_docker_fixtures.sh
```

It's also possible to target mocks against this server, or any other PowerDNS instance which is running on `http://localhost:8080`.

```bash
make test-without-mocks
```

The mocks assume there is a vHost/Server ID called `localhost` and API key `apipw`.
