# PowerDNS 4.1 API bindings for Golang

This community project provides bindings for the currently latest version of PowerDNS.

## Requirements

- PowerDNS 4.1 ("API v1")
  - `--webserver=yes --api=yes --api-key=apipw --api-readonly=no`
  - Note that API v1 is actively maintained. There are differences between 3.x, 4.0 and 4.1 and this client works only with 4.1.
- Go 1.10 (should work with other minor releases as well)

## Installation

```bash
go get github.com/joeig/go-powerdns
```

## Usage

### Initialize the handle

```go
import "github.com/joeig/go-powerdns"

pdns := powerdns.NewClient("http://localhost:80", "localhost", "example.com", "apipw")
```

Assuming that the server is listening on http://localhost:80 for virtual host `localhost`, the API password is `apipw` and you want to edit the domain `example.com`.

### Add/change/delete resource records

```go
zone, err := pdns.AddRecord("www.example.com", "AAAA", 60, ["::1"])
zone, err := pdns.ChangeRecord("www.example.com", "AAAA", 3600, ["::1"])
zone, err := pdns.DeleteRecord("www.example.com", "A")
notifyResult, err := pdns.Notify()
```

### Request zone data

```go
zone, err := pdns.GetZone()
zones, err := pdns.GetZones()
```

### Request server information and statistics

```go
statistics, err := pdns.GetStatistics()
server, err := pdns.GetServer()
servers, err := pdns.GetServers()
```

## Documentation

See [GoDoc](https://godoc.org/github.com/joeig/go-powerdns).

## Contribution

This API client has not been completed yet, so feel free to contribute.

Based on the work of [jgreat](https://github.com/jgreat/powerdns) and [waynz0r](https://github.com/waynz0r/go-powerdns).
