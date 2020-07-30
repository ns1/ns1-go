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

DNS views and compatibility
===========================

DNS views is a means for NS1 to serve one set of data to one group of clients
(e.g. internal employees), and different sets of data to other groups of
clients (e.g. public internet). This has been largely exposed by allowing zones
in the NS1 system to share the same FQDN, and allowing propagation to be
controlled via ACLs and "views". For existing zones, and users that have no
need for the added complexity of views, the default behavior is unchanged.
However, it is important to understand that the requirement that an FQDN be
unique within the network is removed in v3.x, and the ramifications of that.

Views are available in API v3.x. For more information, see:
https://help.ns1.com/hc/en-us/articles/360054071374

COMPATIBILITY

Since an FQDN can now appear in more than one "zone", it can no longer uniquely
identify a zone. Instead, a user-supplied "name", unique within an account,
is used to uniquely identify the zone.

The zone name can be the same as its FQDN, existing zones transfer to the new
schema this way. And if not using views, it does serve as a good identifier.
However, if a second zone is created pointing to same FQDN, it cannot reuse
the FQDN as an identifier, and queries for the zone FQDN (as an identfier)
will return the first zone. The following example should help illustrate:

In general, for API requests, the identifier for a zone will appear in the URL.
They can also be passed or received as fields on an object. In the following
request to a v2.x system, we use ZONE as an identifier in the URL, and may
pass the `zone` field as an indicator of the FQDN for the zone.

$ENDPOINT/v1/zones/example.com -d '{zone: example.com}'
                   ^                ^
                   L___ ZONE (id)   L____ ZONE (fqdn)

Note also that it has been an ERROR if the values in the URL and `zone` field
do not match, and the API would reject.

Going forward, the unique identifier and FQDN of the zone are decoupled. The
reference in the URL is user-assignable, and is passed and returned using the
`name` field for zones (and the `zone_name` field for records). In the new
paradigm, the previous call would look more like:

$ENDPOINT/v1/zones/example.com -d '{zone: example.com, name: example.com}'
                   ^                ^                  ^
                   L___ NAME        L____ ZONE         L____ NAME

For compatibility, if `name` isn't present, the API will use the FQDN, so the
2.x call above should continue to work, and have the same result.

However, in 3.x we are now allowed to make new zones with the same FQDN:

$ENDPOINT/v1/zones/example-internal -d '{zone: example.com, name: example-internal}'
                   ^                     ^                  ^
                   L___ NAME             L____ ZONE         L____ NAME

`example-internal` shares the FQDN with `example.com`. API calls using
`example.com` as the identifier will uniquely identify the first zone, to
address the second zone, `example-internal` must be used in the identifier

So, you can continue to use the FQDN as an identifier - in fact, if you
are not using "views", it is recommended that you do so, but other zones using
the same FQDN will have to choose different names.

Note also that both the example-internal and example.com "zones" can coexist.
How they are propagated relies on how they are organized with regard to views,
acls, and networks.

SUMMARY OF CHANGES TO THE SDK FOR VIEWS

For compatibility, we've left the `zone` variable alone in existing functions.
When name and FQDN differ, and you are calling a func that takes `zone`, you
want to pass the identifier, `name`. The only time you need to be sure to pass
the zone's FQDN is on zone or record creation, where we need to know it.
Otherwise, we're doing lookups, and the `name` is the unambigious field.

For those using views, the NewNamedZone and NewNamedRecord funcs are provided
to create zones and records. If not using views, you can opt to continue using
the older functions.

Contributing
============
Pull Requests and issues are welcome. See the [NS1 Contribution Guidelines](https://github.com/ns1/community) for more information.

Run tests:

```
make test
```

Local dev: use `go mod replace` in client code to point to local checkout of
this repository.

Consider running `./script/install-git-hooks` to install local git hooks for this
project.

# LICENSE

Apache2 - see the included LICENSE file for more information
