package dns

import (
	"bytes"
	"encoding/json"
	"testing"
)

var marshalRecordCases = []struct {
	name    string
	record  *Record
	answers []*Answer
	out     []byte
}{
	{
		"marshalCAARecord",
		NewRecord("example.com", "caa.example.com", "CAA"),
		[]*Answer{NewCAAAnswer(0, "issue", "letsencrypt.org")},
		[]byte(`{"meta":{},"zone_name":"example.com","zone":"example.com","domain":"caa.example.com","type":"CAA","answers":[{"meta":{},"answer":["0","issue","letsencrypt.org"]}],"filters":[]}`),
	},
	{
		"marshalURLFWDRecord",
		NewRecord("example.com", "fwd.example.com", "URLFWD"),
		[]*Answer{
			NewURLFWDAnswer("/net", "https://example.net", 301, 1, 1),
			NewURLFWDAnswer("/org", "https://example.org", 302, 2, 0),
		},
		[]byte(`{"answers":[{"answer":["/net","https://example.net",301,1,1],"meta":{}},{"answer":["/org","https://example.org",302,2,0],"meta":{}}],"meta":{},"zone_name":"example.com","zone":"example.com","domain":"fwd.example.com","type":"URLFWD","filters":[]}`),
	},
}

func TestMarshalRecords(t *testing.T) {
	for _, tt := range marshalRecordCases {
		t.Run(tt.name, func(t *testing.T) {
			for i := range tt.answers {
				tt.record.AddAnswer(tt.answers[i])
			}
			result, err := json.Marshal(tt.record)
			if err != nil {
				t.Error(err)
			}
			if bytes.Compare(result, tt.out) != 0 {
				t.Errorf("got %q, want %q", result, tt.out)
			}
		})
	}
}
