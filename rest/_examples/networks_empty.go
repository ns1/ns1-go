package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dns"
)

func main() {
	fmt.Println("NS1 SDK Empty Networks Example")
	fmt.Println("------------------------------")

	// Case 1: Default behavior (nil Networks)
	zone1 := dns.NewZone("example.com")
	data1, _ := json.Marshal(zone1)
	fmt.Println("Case 1 - Default (nil Networks):")
	fmt.Printf("  JSON: %s\n", string(data1))
	fmt.Printf("  Contains 'networks': %v\n\n", strings.Contains(string(data1), "networks"))

	// Case 2: Empty Networks (explicit empty array)
	zone2 := dns.NewZone("empty.example.com")
	empty := []int{}
	zone2.SetNetworks(&empty)
	data2, _ := json.Marshal(zone2)
	fmt.Println("Case 2 - Empty Networks (explicit empty array):")
	fmt.Printf("  JSON: %s\n", string(data2))
	fmt.Printf("  Contains 'networks': %v\n", strings.Contains(string(data2), "networks"))
	fmt.Printf("  Has empty networks array: %v\n\n", strings.Contains(string(data2), `"networks":[]`))

	// Case 3: Populated Networks
	zone3 := dns.NewZone("populated.example.com")
	networks := []int{1, 2}
	zone3.SetNetworks(&networks)
	data3, _ := json.Marshal(zone3)
	fmt.Println("Case 3 - Populated Networks:")
	fmt.Printf("  JSON: %s\n", string(data3))
	fmt.Printf("  Contains networks: %v\n\n", strings.Contains(string(data3), `"networks":[1,2]`))

	// Case 4: Legacy NetworkIDs with EnsureNetworksFromLegacy
	zone4 := dns.NewZone("legacy.example.com")
	zone4.NetworkIDs = []int{3, 4}
	zone4.EnsureNetworksFromLegacy()
	data4, _ := json.Marshal(zone4)
	fmt.Println("Case 4 - Legacy NetworkIDs with EnsureNetworksFromLegacy:")
	fmt.Printf("  JSON: %s\n", string(data4))
	fmt.Printf("  Contains networks: %v\n", strings.Contains(string(data4), `"networks":[3,4]`))
}
