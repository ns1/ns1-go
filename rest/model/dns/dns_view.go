package dns

// DNSView wraps an NS1 views/ resource
type DNSView struct {
	Name        string   `json:"name,omitempty"`
	Created_at  int      `json:"created_at,omitempty"`
	Updated_at  int      `json:"updated_at,omitempty"`
	Read_acls   []string `json:"read_acls,omitempty"`
	Update_acls []string `json:"update,omitempty"`
	Zones       []string `json:"zones,omitempty"`
	Networks    []int    `json:"networks,omitempty"`
	Preference  int      `json:"preference,omitempty"`
}

// NewDNSView takes a view_name and creates a *DNSView
func NewDNSView(view_name string) *DNSView {
	return &DNSView{
		Name: view_name,
	}
}
