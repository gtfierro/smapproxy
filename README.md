# Smap Proxy

Simple proxy that gives a basic auth HTTP interface to sMAP archiver queries; allows you to put the archiver
listening on localhost and then place this proxy over the top to provide semi-protected access.

## Installation

```bash
git clone github.com/gtfierro/smapproxy
go build
```

## Usage

```bash
$ ./smapproxy -h
Usage of ./smapproxy:
  -archiver string
    	Query URL to proxy (default "http://localhost:8079/api/query")
  -p string
    	Port to serve (default "1212")
  -pass string
    	HTTP basic auth pass
  -user string
    	HTTP basic auth user
$ ./smapproxy -p 8080 -pass username -pass mypassword
```
