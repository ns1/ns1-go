# NS1 Golang SDK

The golang client for the NS1 API: https://api.nsone.net/

# Installing

Just checkout this library to your GOPATH as usual. If you're writing a standard go program
using this library, that should be as simple as saying 'go get'

# Reference documentation

See [the godoc](http://www.godoc.org/github.com/ns1/ns1-go)

# Supported features

## Setup zones
    * Links supported
    * Secondary zones supported

## Setup records in those zones
    * A, MX and CNAME records are supported.
    * Allows records to be linked to other records
    * Allows multiple answers, each of which can be linked to a data feed

## Data sources
    * Can create datasources with arbitrary config

## Data feeds
    * Create data feeds linked to a data source with a label

## NS1 monitoring
    * CRUD monitoring jobs

# Unsupported features

## Users / Account management / API keys
  * Limited

## Useage / querying APIs
  * No support


Examples
========

See [rest example dir]()

Contributing
============

Contributions, ideas and criticisms are all welcome.

# LICENSE

Apache2 - see the included LICENSE file for more information

