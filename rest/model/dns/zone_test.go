package dns

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/rest/model/data"
)

func TestUnmarshalZoneRecords(t *testing.T) {
	d := []byte(`[
    {
      "domain": "foo.test.zone",
      "short_answers": [
        "1.2.3.4"
      ],
      "link": null,
      "ttl": 180,
      "tier": 1.0,
      "type": "A",
      "id": "835d91f01c56932156905bf7"
    },
    {
      "domain": "bar.test.zone",
      "short_answers": [
        "5.6.7.8"
      ],
      "link": null,
      "ttl": 180,
      "tier": 1,
      "type": "A",
      "id": "367d91f01c56932156905v98"
    }
]
`)
	zrl := []*ZoneRecord{}
	if err := json.Unmarshal(d, &zrl); err != nil {
		t.Error(err)
	}

	if len(zrl) != 2 {
		fmt.Println(zrl)
		t.Error("Do not have 2 records in list")
	}

}

func TestUnmarshalZones(t *testing.T) {
	d := []byte(`[
   {
      "nx_ttl":3600,
      "retry":7200,
      "zone":"test.zone",
      "network_pools":[
         "p09"
      ],
      "primary":{
         "enabled":true,
         "secondaries":[
            {
               "ip":"1.1.1.1",
               "notify":true,
               "networks":[

               ],
               "port":53
            },
            {
               "ip":"2.2.2.2",
               "notify":true,
               "networks":[

               ],
               "port":53
            }
         ]
      },
      "refresh":43200,
      "expiry":1209600,
      "dns_servers":[
         "dns1.p09.nsone.net",
         "dns2.p09.nsone.net",
         "dns3.p09.nsone.net",
         "dns4.p09.nsone.net"
      ],
      "meta":{

      },
      "link":null,
      "serial":1473863358,
      "ttl":3600,
      "id":"57d95da659272400013334de",
      "hostmaster":"hostmaster@nsone.net",
      "networks":[
         0
      ],
      "pool":"p09"
   },
   {
     "nx_ttl":3600,
     "retry":7200,
     "zone":"secondary.zone",
     "network_pools":[
        "p09"
     ],
     "secondary":{
        "status":"pending",
        "last_xfr":0,
        "primary_ip":"1.1.1.1",
        "primary_port":53,
      "other_ips":[
      "1.1.1.2",
      "1.1.1.3"
      ],
      "other_ports":[
      53,
      53
      ],
        "enabled":true,
        "tsig":{
     "enabled":false,
     "hash":null,
     "name":null,
     "key":null
        },
        "error":null,
        "expired":false
     },
     "primary":{
        "enabled":false,
        "secondaries":[

        ]
     },
     "refresh":43200,
     "expiry":1209600,
     "dns_servers":[
        "dns1.p09.nsone.net",
        "dns2.p09.nsone.net",
        "dns3.p09.nsone.net",
        "dns4.p09.nsone.net"
     ],
     "meta":{

     },
     "link":null,
     "serial":1473868413,
     "ttl":3600,
     "id":"57d9727d1c372700011eff6e",
     "hostmaster":"hostmaster@nsone.net",
     "networks":[
        0
     ],
     "pool":"p09",
     "dnssec":true
  },
   {
      "nx_ttl":3600,
      "retry":7200,
      "zone":"myfailover.com",
      "network_pools":[
         "p09"
      ],
      "primary":{
         "enabled":true,
         "secondaries":[

         ]
      },
      "refresh":43200,
      "expiry":1209600,
      "dns_servers":[
         "dns1.p09.nsone.net",
         "dns2.p09.nsone.net",
         "dns3.p09.nsone.net",
         "dns4.p09.nsone.net"
      ],
      "meta":{

      },
      "link":null,
      "serial":1473813629,
      "ttl":3600,
      "id":"57d89c7b1c372700011e0a97",
      "hostmaster":"hostmaster@nsone.net",
      "networks":[
         0
      ],
      "pool":"p09",
      "dnssec":false
   }
]
`)
	zl := []*Zone{}
	if err := json.Unmarshal(d, &zl); err != nil {
		t.Error(err)
	}
	if len(zl) != 3 {
		fmt.Println(zl)
		t.Error("Do not have 3 zones in list")
	}
	z := zl[0]
	assert.Nil(t, z.Link)
	assert.Nil(t, z.Secondary)
	assert.Nil(t, z.DNSSEC, "Zone DNSSEC should be nil")
	assert.Equal(t, z.ID, "57d95da659272400013334de", "Wrong zone id")
	assert.Equal(t, z.Zone, "test.zone", "Wrong zone name")
	assert.Equal(t, z.TTL, 3600, "Wrong zone ttl")
	assert.Equal(t, z.NxTTL, 3600, "Wrong zone nxttl")
	assert.Equal(t, z.Retry, 7200, "Wrong zone retry")
	assert.Equal(t, z.Serial, 1473863358, "Wrong zone serial")
	assert.Equal(t, z.Refresh, 43200, "Wrong zone refresh")
	assert.Equal(t, z.Expiry, 1209600, "Wrong zone expiry")
	assert.Equal(t, z.Hostmaster, "hostmaster@nsone.net", "Wrong zone hostmaster")
	assert.Equal(t, z.Pool, "p09", "Wrong zone pool")
	assert.Equal(t, z.NetworkIDs, []int{0}, "Wrong zone network")
	assert.Equal(t, z.NetworkPools, []string{"p09"}, "Wrong zone network pools")
	assert.Equal(t, z.Meta, &data.Meta{}, "Zone meta should be empty")

	dnsServers := []string{
		"dns1.p09.nsone.net",
		"dns2.p09.nsone.net",
		"dns3.p09.nsone.net",
		"dns4.p09.nsone.net",
	}
	assert.Equal(t, z.DNSServers, dnsServers, "Wrong zone dns networks")

	primary := &ZonePrimary{
		Enabled: true,
		Secondaries: []ZoneSecondaryServer{
			ZoneSecondaryServer{
				IP:         "1.1.1.1",
				Port:       53,
				Notify:     true,
				NetworkIDs: []int{},
			},
			ZoneSecondaryServer{
				IP:         "2.2.2.2",
				Port:       53,
				Notify:     true,
				NetworkIDs: []int{},
			},
		},
	}
	assert.Equal(t, z.Primary, primary, "Wrong zone primary")

	// Check zone with secondaries and DNSSEC
	secZ := zl[1]
	assert.Nil(t, secZ.Link)
	assert.Equal(t, *secZ.DNSSEC, true, "Zone DNSSEC should be true")
	assert.Equal(t, secZ.Zone, "secondary.zone", "Wrong zone name")
	assert.Equal(t, secZ.Primary, &ZonePrimary{
		Enabled:     false,
		Secondaries: []ZoneSecondaryServer{}}, "Wrong zone secondary primary")

	secondary := secZ.Secondary
	assert.Nil(t, secondary.Error)
	assert.Equal(t, secondary.Status, "pending", "Wrong zone secondary status")
	assert.Equal(t, secondary.LastXfr, 0, "Wrong zone secondary last xfr")
	assert.Equal(t, secondary.PrimaryIP, "1.1.1.1", "Wrong zone secondary primary ip")
	assert.Equal(t, secondary.PrimaryPort, 53, "Wrong zone secondary primary port")
	assert.Equal(t, secondary.Enabled, true, "Wrong zone secondary enabled")
	assert.ElementsMatch(t, secondary.OtherIPs, []string{"1.1.1.2", "1.1.1.3"}, "Wrong zone secondary list of other IPs")
	assert.ElementsMatch(t, secondary.OtherPorts, []int{53, 53}, "Wrong zone secondary list of other ports")

	assert.Equal(t, secondary.TSIG, &TSIG{
		Enabled: false,
		Hash:    "",
		Name:    "",
		Key:     ""}, "Wrong zone secondary tsig")

	// check last zone with DNSSEC explicitly false
	failoverZ := zl[2]
	assert.Equal(t, *failoverZ.DNSSEC, false, "Zone DNSSEC should be false")
}

func TestMakeSecondary(t *testing.T) {
	d := []byte(`
  {
    "nx_ttl":3600,
    "retry":7200,
    "zone":"test.zone",
    "network_pools":[
       "p09"
    ],
    "primary":{
       "enabled":true,
       "secondaries":[
        {
         "ip":"1.1.1.1",
         "notify":true,
         "networks":[

         ],
         "port":53
        },
        {
         "ip":"2.2.2.2",
         "notify":true,
         "networks":[

         ],
         "port":53
        }
       ]
    },
    "refresh":43200,
    "expiry":1209600,
    "dns_servers":[
       "dns1.p09.nsone.net",
       "dns2.p09.nsone.net",
       "dns3.p09.nsone.net",
       "dns4.p09.nsone.net"
    ],
    "meta":{

    },
    "link":null,
    "serial":1473863358,
    "ttl":3600,
    "id":"57d95da659272400013334de",
    "hostmaster":"hostmaster@nsone.net",
    "networks":[
       0
    ],
    "pool":"p09"
   }
  `)
	z := &Zone{}
	if err := json.Unmarshal(d, &z); err != nil {
		t.Error(err)
	}

	z.MakeSecondary("1.1.1.1")

	primary := &ZonePrimary{
		Enabled:     false,
		Secondaries: make([]ZoneSecondaryServer, 0),
	}
	assert.Equal(t, z.Primary, primary, "Zone primary should be disabled")
	assert.Equal(t, z.Secondary.PrimaryIP, "1.1.1.1", "Wrong zone secondary primary IP")
	assert.Equal(t, z.Secondary.PrimaryPort, 53, "Wrong zone secondary primary port")
}

func TestMarshalSecondaryZone(t *testing.T) {
	z := NewZone("secondary.zone")
	z.MakeSecondary("1.1.1.1")
	z.Secondary.OtherIPs = []string{"2.2.2.2", "3.3.3.3"}
	z.Secondary.OtherPorts = []int{53, 53}
	z.Secondary.TSIG = &TSIG{Enabled: false}

	expected := `
	{
	  "zone": "secondary.zone",
	  "primary": {
		"enabled": false,
		"secondaries": []
	  },
	  "secondary": {
		"error": null,
		"primary_ip": "1.1.1.1",
		"primary_port": 53,
		"other_ips": [
		  "2.2.2.2",
		  "3.3.3.3"
		],
		"other_ports": [
		  53,
		  53
		],
		"enabled": true,
		"tsig": {
		  "enabled": false
		}
	  }
	}`

	j, err := json.Marshal(z)
	require.NoError(t, err)
	assert.JSONEq(t, expected, string(j))
}
