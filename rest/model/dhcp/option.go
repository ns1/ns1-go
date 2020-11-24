package dhcp

// Option encapsulates DHCPv4 and DHCPv6 option information.
//
// Serialized:
//      '{"name": "dhcpv4/routers", "value": ["127.0.0.1"]}'
//      '{"name": "dhcpv4/boot-file-name", "value": "/bootfilename"}'
//      '{"name": "dhcpv6/dns-servers", "value": ["2001:db8::cafe"]}'
type Option struct {
	Name        string      `json:"name"`
	Value       interface{} `json:"value"`
	AlwaysSend  *bool       `json:"always_send,omitempty"`
	Encapsulate *string     `json:"encapsulate,omitempty"`
}

// OptionSet is a convenience type for marshalling an array of options to and from a JSON field.
type OptionSet []Option
