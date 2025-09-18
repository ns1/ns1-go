package dns

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZoneNetworks_MarshalJSON(t *testing.T) {
	// Test Case 1: nil NetworkIDs should be omitted
	z := Zone{
		Zone:       "example.com",
		NetworkIDs: nil,
	}

	data, err := json.Marshal(z)
	assert.NoError(t, err)
	assert.NotContains(t, string(data), "networks")

	// Test Case 2: Empty NetworkIDs should be included as []
	z = Zone{
		Zone:       "example.com",
		NetworkIDs: []int{},
	}

	data, err = json.Marshal(z)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"networks":[]`)

	// Test Case 3: Populated NetworkIDs should be included
	z = Zone{
		Zone:       "example.com",
		NetworkIDs: []int{1, 2, 3},
	}

	data, err = json.Marshal(z)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"networks":[1,2,3]`)
}

func TestZoneNetworks_UnmarshalJSON(t *testing.T) {
	// Test Case 1: JSON with networks field should populate NetworkIDs
	jsonStr := `{"zone":"example.com","networks":[1,2,3]}`
	var z Zone
	err := json.Unmarshal([]byte(jsonStr), &z)
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, z.NetworkIDs)

	// Test Case 2: JSON with empty networks array should result in empty NetworkIDs
	jsonStr = `{"zone":"example.com","networks":[]}`
	z = Zone{}
	err = json.Unmarshal([]byte(jsonStr), &z)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(z.NetworkIDs))

	// Test Case 3: JSON without networks field should result in nil NetworkIDs
	jsonStr = `{"zone":"example.com"}`
	z = Zone{}
	err = json.Unmarshal([]byte(jsonStr), &z)
	assert.NoError(t, err)
	assert.Nil(t, z.NetworkIDs)
}

func TestZone_EnsureNetworksFromLegacy(t *testing.T) {
	// This method is now deprecated but we still test it for compatibility

	// Test Case 1: When networks is nil and NetworkIDs has values
	networkIDs := []int{1, 2, 3}
	z := Zone{
		Zone:       "example.com",
		NetworkIDs: networkIDs,
	}
	// Force the internal networks field to nil
	z.networks = nil

	z.EnsureNetworksFromLegacy()
	// Check that the internal field is now set
	assert.NotNil(t, z.networks)
	assert.Equal(t, networkIDs, *z.networks)

	// Test Case 2: When networks is already set, it shouldn't change
	otherNetworks := []int{4, 5, 6}
	z = Zone{
		Zone:       "example.com",
		NetworkIDs: networkIDs,
	}
	z.networks = &otherNetworks

	z.EnsureNetworksFromLegacy()
	assert.Equal(t, otherNetworks, *z.networks)
	assert.NotEqual(t, networkIDs, *z.networks)

	// Test Case 3: When both are empty/nil
	z = Zone{
		Zone:       "example.com",
		NetworkIDs: nil,
	}
	z.networks = nil

	z.EnsureNetworksFromLegacy()
	assert.Nil(t, z.networks)
	assert.Nil(t, z.NetworkIDs)
}

// Integration test to verify the entire flow works as expected
func TestZoneNetworks_IntegrationFlow(t *testing.T) {
	// Starting with a Zone using NetworkIDs
	z := Zone{
		Zone:       "example.com",
		NetworkIDs: []int{1, 2, 3},
	}

	// Step 1: Marshal to JSON - should automatically populate networks field
	data, err := json.Marshal(z)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"networks":[1,2,3]`)

	// Step 2: Change to empty networks array
	z.NetworkIDs = []int{}

	// Step 3: Marshal again to verify empty array is sent
	data, err = json.Marshal(z)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"networks":[]`)

	// Step 4: Unmarshal from JSON with networks field
	jsonStr := `{"zone":"example.com","networks":[4,5]}`
	err = json.Unmarshal([]byte(jsonStr), &z)
	assert.NoError(t, err)
	assert.Equal(t, []int{4, 5}, z.NetworkIDs)

	// Step 5: Unmarshal from JSON without networks field
	jsonStr = `{"zone":"example.com"}`
	err = json.Unmarshal([]byte(jsonStr), &z)
	assert.NoError(t, err)
	assert.Nil(t, z.NetworkIDs)
}
