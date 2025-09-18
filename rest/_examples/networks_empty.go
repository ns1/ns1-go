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

	// Case 1: Default behavior (nil NetworkIDs)
	zone1 := dns.NewZone("example.com")
	// NetworkIDs is nil by default
	data1, _ := json.Marshal(zone1)
	fmt.Println("Case 1 - Default (nil NetworkIDs):")
	fmt.Printf("  JSON: %s\n", string(data1))
	fmt.Printf("  Contains 'networks': %v\n\n", strings.Contains(string(data1), "networks"))

	// Case 2: Empty NetworkIDs (explicit empty array)
	zone2 := dns.NewZone("empty.example.com")
	zone2.NetworkIDs = []int{} // Empty slice
	data2, _ := json.Marshal(zone2)
	fmt.Println("Case 2 - Empty NetworkIDs (explicit empty array):")
	fmt.Printf("  JSON: %s\n", string(data2))
	fmt.Printf("  Contains 'networks': %v\n", strings.Contains(string(data2), "networks"))
	fmt.Printf("  Has empty networks array: %v\n\n", strings.Contains(string(data2), `"networks":[]`))

	// Case 3: Populated NetworkIDs
	zone3 := dns.NewZone("populated.example.com")
	zone3.NetworkIDs = []int{1, 2}
	data3, _ := json.Marshal(zone3)
	fmt.Println("Case 3 - Populated NetworkIDs:")
	fmt.Printf("  JSON: %s\n", string(data3))
	fmt.Printf("  Contains networks: %v\n\n", strings.Contains(string(data3), `"networks":[1,2]`))

	// Case 4: Unmarshal JSON with networks field
	jsonStr := `{"zone":"example.com","networks":[3,4]}`
	var zone4 dns.Zone
	json.Unmarshal([]byte(jsonStr), &zone4)
	fmt.Println("Case 4 - Unmarshal JSON with networks field:")
	fmt.Printf("  Zone: %s\n", zone4.Zone)
	fmt.Printf("  NetworkIDs: %v\n", zone4.NetworkIDs)
}
