# goproxy

A simple http proxy server, supporting http proxy relay.

### Install

```bash

#正向代理
go get github.com/fondoger/goproxy/forwardproxy@latest

#反向代理
go install github.com/fondoger/goproxy/reverseproxy@latest
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
