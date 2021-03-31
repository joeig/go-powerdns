# PowerDNS 4.x API bindings for Golang

This community project provides bindings for the currently latest version of PowerDNS Authoritative Server.

[![Build Status](https://github.com/joeig/go-powerdns/workflows/Tests/badge.svg)](https://github.com/joeig/go-powerdns/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/joeig/go-powerdns)](https://goreportcard.com/report/github.com/joeig/go-powerdns)
[![Coverage Status](https://coveralls.io/repos/github/joeig/go-powerdns/badge.svg?branch=master)](https://coveralls.io/github/joeig/go-powerdns?branch=master)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/joeig/go-powerdns/v2)](https://pkg.go.dev/github.com/joeig/go-powerdns/v2)

## Features

* Zone handling
* Resource record handling
* Server statistics gathering
* DNSSEC handling

For more features, consult our [documentation](https://pkg.go.dev/github.com/joeig/go-powerdns/v2).

## Usage

### Initialize the handle

```go
import "github.com/joeig/go-powerdns/v2"

pdns := powerdns.NewClient("http://localhost:80", "localhost", map[string]string{"X-API-Key": "apipw"}, nil)
```

Assuming that the server is listening on http://localhost:80 for virtual host `localhost`, the API password is `apipw` and you want to edit the domain `example.com`.

### Get/add/change/delete zones

```go
zones, err := pdns.Zones.List()
zone, err := pdns.Zones.Get("example.com")
export, err := pdns.Zones.Export("example.com")
zone, err := pdns.Zones.AddNative("example.com", true, "", false, "foo", "foo", true, []string{"ns.foo.tld."})
zone, err := pdns.Zones.Add(&powerdns.Zone{})
err := pdns.Zones.Change("example.com", &zone)
err := pdns.Zones.Delete("example.com")
```

### Add/change/delete resource records

```go
err := pdns.Records.Add("example.com", "www.example.com", powerdns.RRTypeAAAA, 60, []string{"::1"})
err := pdns.Records.Change("example.com", "www.example.com", powerdns.RRTypeAAAA, 3600, []string{"::1"})
err := pdns.Records.Delete("example.com", "www.example.com", powerdns.RRTypeA)
```

### Request server information and statistics

```go
statistics, err := pdns.Statistics.List()
servers, err := pdns.Servers.List()
server, err := pdns.Servers.Get("localhost")
```

### Handle DNSSEC cryptographic material

```go
cryptokeys, err := pdns.Cryptokeys.List()
cryptokey, err := pdns.Cryptokeys.Get("example.com", "1337")
err := pdns.Cryptokeys.Delete("example.com", "1337")
```

### More examples

See [examples](https://github.com/joeig/go-powerdns/tree/master/examples).

## Setup

### Requirements

- PowerDNS 4.x ("API v1")
  - `--webserver=yes --api=yes --api-key=apipw --api-readonly=no`
  - Note that API v1 is actively maintained. There are major differences between 3.x, 4.0 and 4.1 and this client works only with 4.1, 4.2 and 4.3.
- Tested with Go version 1.14/1.15/1.16, according to [Go's version support policy](https://golang.org/doc/devel/release.html#policy) (should work with other minor releases as well)

### Install from source

```bash
go get -u github.com/joeig/go-powerdns
```

## Documentation

See [GoDoc](https://godoc.org/github.com/joeig/go-powerdns).

## Contribution

This API client has not been completed yet, so feel free to contribute. The [OpenAPI specification](https://github.com/PowerDNS/pdns/blob/master/docs/http-api/swagger/authoritative-api-swagger.yaml) might be a good reference.

Start a PowerDNS authoritative server including a generic MySQL backend, DNSSEC support and some fixtures using Docker compose:

```bash
docker-compose -f docker-compose-v4.3.yml up
docker-compose -f docker-compose-v4.3.yml exec powerdns sh init_docker_fixtures.sh
```

It's also possible to target mocks against this server:

```bash
make test-without-mocks
```
