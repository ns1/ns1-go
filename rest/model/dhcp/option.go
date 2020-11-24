package dhcp

// Option encapsulates DHCPv{4,6} option information for JSON API and persistence. It pairs an namespaced name with a
// corresponding value; the `Name` should be used to interpret (type assert) the value.
//
// Valid() returns an error if the backing DHCP servers would not accept the value.
//
// Persist example:
//      if err := opt.Valid(); err != nil {
//           return err
//      } else if bytes, err := json.Marshal(opt); err != nil {
//           return err
//      } else {
//           return pushJSONB(db, bytes)
//      }
//
// Value access example:
//      if opt.Name == dhcp.StdOpts.V4.Routers.Name() {
//           ips := opt.Value.([]interface{})
//           ...
//      }
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

// OptionSet is a convenience type for marshalling an array of options to and from a JSONB field.
type OptionSet []Option
