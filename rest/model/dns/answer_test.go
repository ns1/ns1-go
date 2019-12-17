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

var newAnswerCases = []struct {
	name string
	in   *Answer
	out  *Answer
}{
	{
		"testCAAAnswer",
		NewCAAAnswer(0, "issue", "ca.test.zone"),
		&Answer{Meta: &data.Meta{}, Rdata: []string{"0", "issue", "ca.test.zone"}},
	},
	{
		"testURLFWDAnswer",
		NewURLFWDAnswer("/", "https://google.com", 302, 2, 0),
		&Answer{
			Meta:  &data.Meta{},
			Rdata: []string{"/", "https://google.com", "302", "2", "0"},
		},
	},
}

func TestNewAnswers(t *testing.T) {
	for _, tt := range newAnswerCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.in, tt.out, tt.name)
		})
	}
}

func TestPrepareURLFWDAnswer(t *testing.T) {
	tests := []struct {
		name   string
		in     *Answer
		out    *Alias
		expErr bool
	}{
		{
			"valid",
			&Answer{Rdata: []string{"/", "https://google.com", "302", "2", "0"}},
			&Alias{
				Rdata:       []interface{}{"/", "https://google.com", 302, 2, 0},
				AliasAnswer: &AliasAnswer{Rdata: []string{"/", "https://google.com", "302", "2", "0"}},
			},
			false,
		},
		{
			"invalid",
			&Answer{Rdata: []string{"/", "https://google.com", "302", "2"}},
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prepareURLFWDAnswer(tt.in)
			if tt.expErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.out, got.(*Alias))
		})
	}
}
