package main

import (
	"fmt"
	"github.com/bobtfish/go-nsone-api"
)

func main() {
	if k := os.Getenv("NSONE_APIKEY"); k == "" {
		fmt.Println("NSONE_APIKEY environment variable is not set, giving up")
	}
	n := nsone.New(k)
	fmt.Println(n.GetZones())
	z := nsone.NewZone("foo.com")
	z.Hostmaster = "hostmaster@foo.com"
	n.CreateZone(z)
	n.DeleteZone("foo.com")
}
