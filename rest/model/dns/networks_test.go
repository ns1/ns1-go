package dns

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZoneNetworks_MarshalJSON(t *testing.T) {
	// Test Case 1: nil Networks should be omitted
	z := Zone{
		Zone:       "example.com",
		Networks:   nil,
		NetworkIDs: nil,
	}
	data, err := json.Marshal(z)
	assert.NoError(t, err)
	assert.NotContains(t, string(data), "networks")

	// Test Case 2: Empty Networks should be included as []
	empty := []int{}
	z = Zone{
		Zone:       "example.com",
		Networks:   &empty,
		NetworkIDs: nil,
	}
	data, err = json.Marshal(z)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"networks":[]`)

	// Test Case 3: Populated Networks should be included
	networks := []int{1, 2, 3}
	z = Zone{
		Zone:       "example.com",
		Networks:   &networks,
		NetworkIDs: nil,
	}
	data, err = json.Marshal(z)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"networks":[1,2,3]`)
}

func TestZoneNetworks_UnmarshalJSON(t *testing.T) {
	// Test Case 1: JSON with networks field should populate both Networks and NetworkIDs
	jsonStr := `{"zone":"example.com","networks":[1,2,3]}`
	var z Zone
	err := json.Unmarshal([]byte(jsonStr), &z)
	assert.NoError(t, err)
	assert.NotNil(t, z.Networks)
	assert.Equal(t, []int{1, 2, 3}, *z.Networks)
	assert.Equal(t, []int{1, 2, 3}, z.NetworkIDs)

	// Test Case 2: JSON with empty networks array should result in empty slices
	jsonStr = `{"zone":"example.com","networks":[]}`
	z = Zone{}
	err = json.Unmarshal([]byte(jsonStr), &z)
	assert.NoError(t, err)
	assert.NotNil(t, z.Networks)
	assert.Equal(t, 0, len(*z.Networks))
	assert.Equal(t, 0, len(z.NetworkIDs))

	// Test Case 3: JSON without networks field should result in nil fields
	jsonStr = `{"zone":"example.com"}`
	z = Zone{}
	err = json.Unmarshal([]byte(jsonStr), &z)
	assert.NoError(t, err)
	assert.Nil(t, z.Networks)
	assert.Nil(t, z.NetworkIDs)
}

func TestZone_EnsureNetworksFromLegacy(t *testing.T) {
	// Test Case 1: When Networks is nil and NetworkIDs has values
	networkIDs := []int{1, 2, 3}
	z := Zone{
		Zone:       "example.com",
		Networks:   nil,
		NetworkIDs: networkIDs,
	}
	z.EnsureNetworksFromLegacy()
	assert.NotNil(t, z.Networks)
	assert.Equal(t, networkIDs, *z.Networks)

	// Test Case 2: When Networks is already set, it shouldn't change
	networks := []int{4, 5, 6}
	z = Zone{
		Zone:       "example.com",
		Networks:   &networks,
		NetworkIDs: networkIDs,
	}
	z.EnsureNetworksFromLegacy()
	assert.NotNil(t, z.Networks)
	assert.Equal(t, networks, *z.Networks)
	assert.NotEqual(t, networkIDs, *z.Networks)

	// Test Case 3: When both are empty/nil
	z = Zone{
		Zone:       "example.com",
		Networks:   nil,
		NetworkIDs: nil,
	}
	z.EnsureNetworksFromLegacy()
	assert.Nil(t, z.Networks)
	assert.Nil(t, z.NetworkIDs)
}

// Integration test to verify the entire flow works as expected
func TestZoneNetworks_IntegrationFlow(t *testing.T) {
	// Starting with a Zone using legacy NetworkIDs
	z := Zone{
		Zone:       "example.com",
		NetworkIDs: []int{1, 2, 3},
	}

	// Step 1: Call EnsureNetworksFromLegacy to populate Networks
	z.EnsureNetworksFromLegacy()
	assert.NotNil(t, z.Networks)
	assert.Equal(t, z.NetworkIDs, *z.Networks)

	// Step 2: Marshal to JSON
	data, err := json.Marshal(z)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"networks":[1,2,3]`)

	// Step 3: Change to empty networks
	empty := []int{}
	z.Networks = &empty
	z.NetworkIDs = nil

	// Step 4: Marshal again to verify empty array is sent
	data, err = json.Marshal(z)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"networks":[]`)

	// Step 5: Unmarshal from JSON with networks field
	jsonStr := `{"zone":"example.com","networks":[4,5]}`
	err = json.Unmarshal([]byte(jsonStr), &z)
	assert.NoError(t, err)
	assert.NotNil(t, z.Networks)
	assert.Equal(t, []int{4, 5}, *z.Networks)
	assert.Equal(t, []int{4, 5}, z.NetworkIDs)
}
