package main

import (
	"fmt"
	"github.com/bobtfish/go-nsone-api"
	"os"
)

func main() {
	k := os.Getenv("NSONE_APIKEY")
	if k == "" {
		fmt.Println("NSONE_APIKEY environment variable is not set, giving up")
	}
	api := nsone.New(k)
	api.Debug()
	api.RateLimitStrategySleep()
	fmt.Println(api.GetZones())

	z := nsone.NewZone("foo.com")
	api.CreateZone(z)

	r := nsone.NewRecord("foo.com", "www.foo.com", "A")
	r.Answers = []nsone.Answer{nsone.Answer{Answer: []string{"1.1.1.1"}}}
	api.CreateRecord(r)

	api.DeleteZone("foo.com")
}
