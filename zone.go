package ns1

// Zone wraps an NS1 /zone resource
type Zone struct {
	ID           string            `json:"id,omitempty"`
	TTL          int               `json:"ttl,omitempty"`
	NxTTL        int               `json:"nx_ttl,omitempty"`
	Retry        int               `json:"retry,omitempty"`
	Zone         string            `json:"zone,omitempty"`
	Refresh      int               `json:"refresh,omitempty"`
	Expiry       int               `json:"expiry,omitempty"`
	Primary      *ZonePrimary      `json:"primary,omitempty"`
	DNSServers   []string          `json:"dns_servers,omitempty"`
	Networks     []int             `json:"networks,omitempty"`
	NetworkPools []string          `json:"network_pools,omitempty"`
	Hostmaster   string            `json:"hostmaster,omitempty"`
	Pool         string            `json:"pool,omitempty"`
	Meta         map[string]string `json:"meta,omitempty"`
	Secondary    *ZoneSecondary    `json:"secondary,omitempty"`
	Link         string            `json:"link,omitempty"`
	Records      []ZoneRecord      `json:"records,omitempty"`
	Serial       int               `json:"serial,omitempty"`
}

// Implementation of Stringer interface.
func (z Zone) String() string {
	return z.Zone
}

// ZoneRecord wraps Zone's "records" attribute
type ZoneRecord struct {
	Domain   string   `json:"Domain,omitempty"`
	ID       string   `json:"id,omitempty"`
	Link     string   `json:"link,omitempty"`
	ShortAns []string `json:"short_answers,omitempty"`
	Tier     int      `json:"tier,omitempty"`
	TTL      int      `json:"ttl,omitempty"`
	Type     string   `json:"type,omitempty"`
}

// ZonePrimary wraps a Zone's "primary" attribute
type ZonePrimary struct {
	Enabled     bool                  `json:"enabled"`
	Secondaries []ZoneSecondaryServer `json:"secondaries"`
}

// ZoneSecondary wraps a Zone's "secondary" attribute
type ZoneSecondary struct {
	Status      string `json:"status,omitempty"`
	LastXfr     int    `json:"last_xfr,omitempty"`
	PrimaryIP   string `json:"primary_ip,omitempty"`
	PrimaryPort int    `json:"primary_port,omitempty"`
	Enabled     bool   `json:"enabled"`
	Expired     bool   `json:"expired,omitempty"`
}

// ZoneSecondaryServer wraps elements of a Zone's "primary.secondary" attribute
type ZoneSecondaryServer struct {
	IP     string `json:"ip"`
	Port   int    `json:"port,omitempty"`
	Notify bool   `json:"notify"`
}

// NewZone takes a zone domain name and creates a new primary *Zone
func NewZone(zone string) *Zone {
	z := Zone{
		Zone: zone,
	}
	z.MakePrimary()
	return &z
}

// MakePrimary enables Primary, disables Secondary, and sets primary's Secondaries to all provided ZoneSecondaryServers
func (z *Zone) MakePrimary(secondaries ...ZoneSecondaryServer) {
	z.Secondary = nil
	z.Primary = &ZonePrimary{
		Enabled:     true,
		Secondaries: secondaries,
	}
	if z.Primary.Secondaries == nil {
		z.Primary.Secondaries = make([]ZoneSecondaryServer, 0)
	}
}

// MakeSecondary enables Secondary, disables Primary, and sets secondary's Primary_ip to provided ip
func (z *Zone) MakeSecondary(ip string) {
	z.Secondary = &ZoneSecondary{
		Enabled:     true,
		PrimaryIP:   ip,
		PrimaryPort: 53,
	}
	z.Primary = &ZonePrimary{
		Enabled:     false,
		Secondaries: make([]ZoneSecondaryServer, 0),
	}
}

// LinkTo sets Link to a target zone domain name and unsets all other configuration properties
func (z *Zone) LinkTo(to string) {
	z.Meta = nil
	z.TTL = 0
	z.NxTTL = 0
	z.Retry = 0
	z.Refresh = 0
	z.Expiry = 0
	z.Primary = nil
	z.DNSServers = nil
	z.Networks = nil
	z.NetworkPools = nil
	z.Hostmaster = ""
	z.Pool = ""
	z.Secondary = nil
	z.Link = to
}
