package filter

// NewSticky returns a filter that sorts answers uniquely depending
// on the IP address of the requester. The same requester always
// gets the same ordering of answers. byNetwork indicates whether
// to apply the 'stickyness' by subnet(not individual IP).
func NewSticky(byNetwork bool) *Filter {
	config := map[string]interface{}{
		"sticky_by_network": byNetwork,
	}
	return &Filter{
		Filter: "sticky",
		Config: config,
	}
}

// NewWeightedSticky returns a filter that shuffles answers randomly
// per-requester based on weight. byNetwork indicates whether to
// apply the 'stickyness' by subnet(not individual IP).
func NewWeightedSticky(byNetwork bool) *Filter {
	config := map[string]interface{}{
		"sticky_by_network": byNetwork,
	}
	return &Filter{
		Filter: "weighted_sticky",
		Config: config,
	}
}

// NewIPv4PrefixShuffle returns a filter that randomly selects
// IPv4 addresses from prefix list. This filter can only be used
// A records. n is the number of IPs to randomly select per answer.
func NewIPv4PrefixShuffle(n int) *Filter {
	config := map[string]interface{}{"N": n}
	return &Filter{
		Filter: "ipv4_prefix_shuffle",
		Config: config,
	}
}

// NewNetfenceASN returns a filter that restricts to answers where
// the ASN of requester IP matches ASN list. rmNoASN determines
// whether to remove answers without asn list on any match.
func NewNetfenceASN(rmNoASN bool) *Filter {
	config := map[string]interface{}{
		"remove_no_asn": rmNoASN,
	}
	return &Filter{
		Filter: "netfence_asn",
		Config: config,
	}
}

// NewNetfencePrefix returns a filter that restricts to answers where
// requester IP matches prefix list. rmNoIPPrefix determines
// whether to remove answers without ip prefixes on any match.
func NewNetfencePrefix(rmNoIPPrefix bool) *Filter {
	config := map[string]interface{}{
		"remove_no_ip_prefixes": rmNoIPPrefix,
	}
	return &Filter{
		Filter: "netfence_prefix",
		Config: config,
	}
}
