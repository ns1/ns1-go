package main

import (
	"fmt"
	"github.com/bobtfish/go-nsone-api"
)

func main() {
	n := nsone.New("fF4uiAIL9wntC7Bdar0a")
	fmt.Println(n.GetZones())
}
