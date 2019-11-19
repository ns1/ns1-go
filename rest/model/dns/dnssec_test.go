package dns

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalDNSSEC(t *testing.T) {
	j := []byte(`{
  "keys": {
    "dnskey": [
     [
        "257",
        "3",
        "13",
        "t+4DPP+MFZ0Cr7gAXiDYv6HTyXzq/O2ESVRLc/ysuh5xBXKIsjsj5baV1HzhBNo2F7mbsevsEo0/6UEL8+JBmA=="
      ],
      [
        "256",
        "3",
        "13",
        "pxEUulkf8UZtE9fy2+4wJwM44xncypgGVps4hE4kQGA5TuC/XJPoKBX6e3B/QL9AmwFCgyFeC4uRNxoqxK0xOg=="
      ]
    ],
    "ttl": 3600
  },
  "delegation": {
    "dnskey": [
      [
        "257",
        "3",
        "13",
        "t+4DPP+MFZ0Cr7gAXiDYv6HTyXzq/O2ESVRLc/ysuh5xBXKIsjsj5baV1HzhBNo2F7mbsevsEo0/6UEL8+JBmA=="
      ]
    ],
    "ds": [
      [
        "48553",
        "13",
        "2",
        "150ae338f365a05e53cb781aedd1b54bf5f27f6a837441292ccf03ca26ad0fb3"
      ]
    ],
    "ttl": 3666
  },
  "zone": "burping.cat."
}
`)
	d := ZoneDNSSEC{}
	if err := json.Unmarshal(j, &d); err != nil {
		t.Error(err)
	}

	assert.Equal(t, d.Zone, "burping.cat.")

	keys := d.Keys
	assert.Equal(t, json.Number("3600"), keys.TTL)
	assert.Equal(t, 2, len(keys.DNSKey))
	k := keys.DNSKey[0]
	assert.Equal(t, "257", k.Flags)
	assert.Equal(t, "3", k.Protocol)
	assert.Equal(t, "13", k.Algorithm)
	assert.Equal(t, "t+4DPP+MFZ0Cr7gAXiDYv6HTyXzq/O2ESVRLc/ysuh5xBXKIsjsj5baV1HzhBNo2F7mbsevsEo0/6UEL8+JBmA==", k.PublicKey)
	k = keys.DNSKey[1]
	assert.Equal(t, "256", k.Flags)
	assert.Equal(t, "3", k.Protocol)
	assert.Equal(t, "13", k.Algorithm)
	assert.Equal(t, "pxEUulkf8UZtE9fy2+4wJwM44xncypgGVps4hE4kQGA5TuC/XJPoKBX6e3B/QL9AmwFCgyFeC4uRNxoqxK0xOg==", k.PublicKey)

	delegation := d.Delegation
	assert.Equal(t, json.Number("3666"), delegation.TTL)
	assert.Equal(t, 1, len(delegation.DNSKey))
	k = delegation.DNSKey[0]
	assert.Equal(t, "257", k.Flags)
	assert.Equal(t, "3", k.Protocol)
	assert.Equal(t, "13", k.Algorithm)
	assert.Equal(t, "t+4DPP+MFZ0Cr7gAXiDYv6HTyXzq/O2ESVRLc/ysuh5xBXKIsjsj5baV1HzhBNo2F7mbsevsEo0/6UEL8+JBmA==", k.PublicKey)
	assert.Equal(t, 1, len(delegation.DS))
	k = delegation.DS[0]
	assert.Equal(t, "48553", k.Flags)
	assert.Equal(t, "13", k.Protocol)
	assert.Equal(t, "2", k.Algorithm)
	assert.Equal(t, "150ae338f365a05e53cb781aedd1b54bf5f27f6a837441292ccf03ca26ad0fb3", k.PublicKey)
}
