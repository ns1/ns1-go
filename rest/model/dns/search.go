package dns

type SearchResult struct {
	Next         string          `json:"next"`
	Limit        int             `json:"limit"`
	TotalResults int             `json:"total_results"`
	Results      []*SearchRecord `json:"results"`
}

type SearchRecord struct {
	Domain     string    `json:"domain"`
	Type       string    `json:"type"`
	TTL        int       `json:"ttl"`
	ZoneFQDN   string    `json:"zone_fqdn"`
	ZoneHandle string    `json:"zone_handle"`
	Answers    []*Answer `json:"answers"`
}

// New correct types
type SearchResultV2 struct {
	Next         string            `json:"next"`
	Limit        int               `json:"limit"`
	TotalResults int               `json:"total_results"`
	Results      []*SearchRecordV2 `json:"results"`
}

type SearchRecordV2 struct {
	Domain     string          `json:"domain"`
	Type       string          `json:"type"`
	TTL        int             `json:"ttl"`
	ZoneFQDN   string          `json:"zone_fqdn"`
	ZoneHandle string          `json:"zone_handle"`
	Answers    []*SearchAnswer `json:"answers"`
}

type SearchAnswer struct {
	Answer []string `json:"answer"`
}
