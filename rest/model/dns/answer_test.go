package dns

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/ns1/ns1-go.v2/rest/model/data"
)

func TestUnmarshalAnswer(t *testing.T) {
	d := []byte(`{
		"id": "520519509f782d58bb4df418",
		"answer": [
			5,
			"record.test.zone"
		],
		"meta": {
			"up": { "feed": "520533b89f782d5b1a10a851" },
			"priority": 1
		},
		"region": "us-east"
	}`)
	a := Answer{}
	if err := json.Unmarshal(d, &a); err != nil {
		t.Error(err)
	}
	assert.Equal(t, "520519509f782d58bb4df418", a.ID, "Incorrect ID")
	assert.Equal(t, []string{"5", "record.test.zone"}, a.Rdata, "Incorrect rdata")
	expectedMeta := &data.Meta{
		Up:       map[string]interface{}{"feed": "520533b89f782d5b1a10a851"},
		Priority: float64(1),
	}
	assert.Equal(t, expectedMeta, a.Meta, "Incorrect meta")
	assert.Equal(t, "us-east", a.RegionName, "Incorrect region")
}

func TestNewCAAAnswer(t *testing.T) {
	expected := &Answer{
		Meta: &data.Meta{},
		Rdata: []string{
			"0",
			"issue",
			"ca.test.zone",
		},
	}
	assert.Equal(t, expected, NewCAAAnswer(0, "issue", "ca.test.zone"))
}
