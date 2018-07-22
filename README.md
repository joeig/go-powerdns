# PowerDNS 4.1 API bindings for Golang

This community project provides bindings for the currently latest version of PowerDNS.

[![Build Status](https://travis-ci.org/joeig/go-powerdns.svg?branch=master)](https://travis-ci.org/joeig/go-powerdns)

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

pdns := powerdns.NewClient("http://localhost:80", "localhost", "apipw")
```

Assuming that the server is listening on http://localhost:80 for virtual host `localhost`, the API password is `apipw` and you want to edit the domain `example.com`.

### Request zone data

```go
zone, err := pdns.GetZone("example.com")
zones, err := pdns.GetZones()
```

### Add/change/delete resource records

```go
err := zone.AddRecord("www.example.com", "AAAA", 60, ["::1"])
err := zone.ChangeRecord("www.example.com", "AAAA", 3600, ["::1"])
err := zone.DeleteRecord("www.example.com", "A")
notifyResult, err := zone.Notify()
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
