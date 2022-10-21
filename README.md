[![Build Status](https://travis-ci.org/ns1/ns1-go.svg?branch=v2)](https://travis-ci.org/ns1/ns1-go) [![GoDoc](https://godoc.org/gopkg.in/ns1/ns1-go.v2?status.svg)](https://godoc.org/gopkg.in/ns1/ns1-go.v2)

# NS1 Golang SDK

> This project is in [active development](https://github.com/ns1/community/blob/master/project_status/ACTIVE_DEVELOPMENT.md).

The golang client for the NS1 API: https://ns1.com/api/

# Installing

```
$ go get gopkg.in/ns1/ns1-go.v2
```

Examples
========

[See more](https://github.com/ns1/ns1-go/tree/v2/rest/_examples)


```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	api "gopkg.in/ns1/ns1-go.v2/rest"
)

func main() {
	k := os.Getenv("NS1_APIKEY")
	if k == "" {
		fmt.Println("NS1_APIKEY environment variable is not set, giving up")
		os.Exit(1)
	}

	httpClient := &http.Client{Timeout: time.Second * 10}
	client := api.NewClient(httpClient, api.SetAPIKey(k))

	zones, _, err := client.Zones.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, z := range zones {
		fmt.Println(z.Zone)
	}

}
```

DNS views and compatibility with pre-2.6.6 SDK versions
=======================================================

DNS views allow NS1 to serve one set of data to one group of clients
(e.g. internal employees), and different results to other groups of
clients (e.g. public internet). Multiple zones can now have the same
fully-qualified domain name (FQDN), with propagation controlled via ACLs
and the Views feature. For more information, please refer to this
[NS1 documentation page](https://help.ns1.com/hc/en-us/articles/360054071374).

Users who do not need views can ignore this feature. Users who do
use views must now use the user-supplied `name` field in the API
to uniquely identify a zone. More than one zone can have the same
FQDN, but their `name` fields must be unique.

For compatibility, the `zone` field is unchanged in existing functions.
When `name` and FQDN differ, and you are calling a func that takes `zone`,
you must add the `name` identifier as well. The zone's FQDN is only
required during zone or record creation.

When using views, the NewNamedZone and NewNamedRecord funcs are provided
to create zones and records. If not using views, you can continue using
the older functions.

Contributing
============
Pull Requests and issues are welcome. See the [NS1 Contribution Guidelines](https://github.com/ns1/community) for more information.

Run tests:

```
make test
```

Local dev: use the `go.mod` `replace` directive or
[`go work use`](https://go.dev/ref/mod#go-work-use)
in client code to point to the local checkout of this repository.

Consider running `./script/install-git-hooks` to install local git hooks for this
project.

# LICENSE

Apache2 - see the included LICENSE file for more information
