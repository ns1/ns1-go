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
