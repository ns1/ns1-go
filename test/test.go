package main

import (
	"fmt"
	"github.com/bobtfish/go-nsone-api"
)

func main() {
	n := nsone.New("fF4uiAIL9wntC7Bdar0a")
	fmt.Println(n.GetZones())
	z := nsone.NewZone("foo.com")
	z.Hostmaster = "hostmaster@foo.com"
	n.CreateZone(z)
	n.DeleteZone(z)
}
