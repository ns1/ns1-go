package ns1

import "fmt"

const (
	zonePath = "zones"
)

// Zone wraps an NS1 /zone resource
type Zone struct {
	Id            string            `json:"id,omitempty"`
	Ttl           int               `json:"ttl,omitempty"`
	Nx_ttl        int               `json:"nx_ttl,omitempty"`
	Retry         int               `json:"retry,omitempty"`
	Zone          string            `json:"zone,omitempty"`
	Refresh       int               `json:"refresh,omitempty"`
	Expiry        int               `json:"expiry,omitempty"`
	Primary       *ZonePrimary      `json:"primary,omitempty"`
	Dns_servers   []string          `json:"dns_servers,omitempty"`
	Networks      []int             `json:"networks,omitempty"`
	Network_pools []string          `json:"network_pools,omitempty"`
	Hostmaster    string            `json:"hostmaster,omitempty"`
	Pool          string            `json:"pool,omitempty"`
	Meta          map[string]string `json:"meta,omitempty"`
	Secondary     *ZoneSecondary    `json:"secondary,omitempty"`
	Link          string            `json:"link,omitempty"`
	Records       []ZoneRecord      `json:"records,omitempty"`
	Serial        int               `json:"serial,omitempty"`
}

// Implementation of Stringer interface.
func (z Zone) String() string {
	return z.Zone
}

// ZoneRecord wraps Zone's "records" attribute
type ZoneRecord struct {
	Domain   string   `json:"Domain,omitempty"`
	Id       string   `json:"id,omitempty"`
	Link     string   `json:"link,omitempty"`
	ShortAns []string `json:"short_answers,omitempty"`
	Tier     int      `json:"tier,omitempty"`
	Ttl      int      `json:"ttl,omitempty"`
	Type     string   `json:"type,omitempty"`
}

// ZonePrimary wraps a Zone's "primary" attribute
type ZonePrimary struct {
	Enabled     bool                  `json:"enabled"`
	Secondaries []ZoneSecondaryServer `json:"secondaries"`
}

// ZoneSecondary wraps a Zone's "secondary" attribute
type ZoneSecondary struct {
	Status       string `json:"status,omitempty"`
	Last_xfr     int    `json:"last_xfr,omitempty"`
	Primary_ip   string `json:"primary_ip,omitempty"`
	Primary_port int    `json:"primary_port,omitempty"`
	Enabled      bool   `json:"enabled"`
	Expired      bool   `json:"expired,omitempty"`
}

// ZoneSecondaryServer wraps elements of a Zone's "primary.secondary" attribute
type ZoneSecondaryServer struct {
	Ip     string `json:"ip"`
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
		Enabled:      true,
		Primary_ip:   ip,
		Primary_port: 53,
	}
	z.Primary = &ZonePrimary{
		Enabled:     false,
		Secondaries: make([]ZoneSecondaryServer, 0),
	}
}

// LinkTo sets Link to a target zone domain name and unsets all other configuration properties
func (z *Zone) LinkTo(to string) {
	z.Meta = nil
	z.Ttl = 0
	z.Nx_ttl = 0
	z.Retry = 0
	z.Refresh = 0
	z.Expiry = 0
	z.Primary = nil
	z.Dns_servers = nil
	z.Networks = nil
	z.Network_pools = nil
	z.Hostmaster = ""
	z.Pool = ""
	z.Secondary = nil
	z.Link = to
}

type ZonesService service

// List returns all active zones and basic zone configuration details for each.
//
// NS1 API docs: https://ns1.com/api/#zones-get
func (s *ZonesService) List() ([]*Zone, error) {
	req, err := s.client.NewRequest("GET", zonePath, nil)
	if err != nil {
		return nil, err
	}

	zl := []*Zone{}
	_, err = s.client.Do(req, &zl)
	if err != nil {
		return nil, err
	}

	return zl, nil
}

// Get takes a zone name and returns a single active zone and its basic configuration details.
//
// NS1 API docs: https://ns1.com/api/#zones-zone-get
func (s *ZonesService) Get(zone string) (*Zone, error) {
	path := fmt.Sprintf("%s/%s", zonePath, zone)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var z Zone
	_, err = s.client.Do(req, &z)
	if err != nil {
		return nil, err
	}

	return &z, nil
}

// Create takes a *Zone and creates a new DNS zone.
//
// NS1 API docs: https://ns1.com/api/#zones-put
func (s *ZonesService) Create(z *Zone) error {
	path := fmt.Sprintf("%s/%s", zonePath, z.Zone)

	req, err := s.client.NewRequest("PUT", path, &z)
	if err != nil {
		return err
	}

	// Update zones fields with data from api(ensure consistent)
	_, err = s.client.Do(req, &z)
	if err != nil {
		return err
	}

	return nil
}

// Update takes a *Zone and modifies basic details of a DNS zone.
//
// NS1 API docs: https://ns1.com/api/#zones-post
func (s *ZonesService) Update(z *Zone) error {
	path := fmt.Sprintf("%s/%s", zonePath, z.Zone)

	req, err := s.client.NewRequest("POST", path, &z)
	if err != nil {
		return err
	}

	// Update zones fields with data from api(ensure consistent)
	_, err = s.client.Do(req, &z)
	if err != nil {
		return err
	}

	return nil
}

// Delete takes a zone and destroys an existing DNS zone and all records in the zone.
//
// NS1 API docs: https://ns1.com/api/#zones-delete
func (s *ZonesService) Delete(zone string) error {
	path := fmt.Sprintf("%s/%s", zonePath, zone)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}
