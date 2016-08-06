package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	ns1 "github.com/ns1/ns1-go"
	"github.com/ns1/ns1-go/rest"
)

type byTTL []*ns1.Zone

func (x byTTL) Len() int           { return len(x) }
func (x byTTL) Less(i, j int) bool { return x[i].TTL < x[j].TTL }
func (x byTTL) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type customSort struct {
	z    []*ns1.Zone
	less func(x, y *ns1.Zone) bool
}

func (x customSort) Len() int           { return len(x.z) }
func (x customSort) Less(i, j int) bool { return x.less(x.z[i], x.z[j]) }
func (x customSort) Swap(i, j int)      { x.z[i], x.z[j] = x.z[j], x.z[i] }

func printZones(zones []*ns1.Zone) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Zone", "TTL", "NxTTL", "Retry", "Refresh")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, z := range zones {
		fmt.Fprintf(tw, format, z.Zone, z.TTL, z.NxTTL, z.Retry, z.Refresh)
	}
	tw.Flush() // calculate column widths and print table
}

func main() {
	k := os.Getenv("NS1_APIKEY")
	if k == "" {
		fmt.Println("NS1_APIKEY environment variable is not set, giving up")
	}

	httpClient := &http.Client{Timeout: time.Second * 10}
	client := rest.NewAPIClient(httpClient, rest.SetAPIKey(k))

	zones, err := client.Zones.List()
	if err != nil {
		log.Fatal(err)
	}

	sort.Sort(byTTL(zones))
	printZones(zones)

	sort.Sort(customSort{zones, func(x, y *ns1.Zone) bool {
		if x.Zone != y.Zone {
			return x.Zone < y.Zone
		}
		if x.Retry != y.Retry {
			return x.Retry < y.Retry
		}
		if x.Refresh != y.Refresh {
			return x.Refresh < y.Refresh
		}
		return false
	}})
	printZones(zones)

}
