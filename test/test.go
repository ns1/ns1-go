package main

import (
	"fmt"
	"os"

	"github.com/ns1/ns1-go"
)

func main() {
	k := os.Getenv("NS1_APIKEY")
	if k == "" {
		fmt.Println("NS1_APIKEY environment variable is not set, giving up")
	}
	api := ns1.New(k)
	api.Debug()
	api.RateLimitStrategySleep()
	fmt.Println(api.GetZones())

	z := ns1.NewZone("foo.com")
	api.CreateZone(z)

	r := ns1.NewRecord("foo.com", "www.foo.com", "A")
	r.Answers = []ns1.Answer{ns1.Answer{Answer: []string{"1.1.1.1"}}}
	api.CreateRecord(r)

	api.DeleteZone("foo.com")
}
