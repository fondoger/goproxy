goproxy
==========

A simple http proxy server, supporting http proxy relay.

### Install
```bash
go get github.com/fondoger/goproxy
```

### Usage
```
Usage: goproxy --addr=0.0.0.0:8080 

Options:
  -addr string
        Listen host port, eg: --addr=0.0.0.0:8000
  -http-relay string
        (optional) Relay to http proxy, eg: --http-relay=http://127.0.0.1:8081
```
