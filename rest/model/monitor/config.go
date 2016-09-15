package monitor

// Config is a flat mapping where values are simple (no slices/maps).
type Config map[string]interface{}

// NewHTTPConfig constructs/returns a job configuration for HTTP type jobs.
// url is the URL to query. (Required)
// method is the HTTP method(valid methods are HEAD, GET, and POST).
// ua is the user agent text in the request header.
// auth is the authorization header to use in request.
// connTimeout is the timeout(in sec) to wait for query output.
func NewHTTPConfig(url, method, ua, auth string, connTimeout int) *Config {
	return &Config{
		"url":                url, // Required
		"method":             method,
		"user_agent":         ua,
		"auth":               auth,
		"connection_timeout": connTimeout,
	}

}

// NewDNSConfig constructs/returns a job configuration for DNS type jobs.
// host is the IP address or hostname of the nameserver to query. (Required)
// domain name to query. (Required)
// port is the dns port to query on host.
// t is the type of the DNS record type to query.
// respTimeout is the timeout(in ms) after sending query to wait for the output.
func NewDNSConfig(host, domain string, port int, t string, respTimeout int) *Config {
	return &Config{
		"host":             host,   // Required
		"domain":           domain, // Required
		"port":             port,
		"type":             t,
		"response_timeout": respTimeout,
	}
}

// NewTCPConfig constructs/returns a job configuration for TCP type jobs.
// host is the IP address or hostname to connect to. (Required)
// port is the tcp port to connect to on host. (Required)
// connTimeout is the timeout(in ms) before giving up on trying to connect.
// respTimeout is the timeout(in sec) after connecting to wait for output.
// send is the string to send to the host upon connecting.
// ssl determines whether to attempt negotiating an SSL connection.
func NewTCPConfig(host string, port, connTimeout, respTimeout int, send string, ssl bool) *Config {
	return &Config{
		"host":               host, // Required
		"port":               port, // Required
		"connection_timeout": connTimeout,
		"response_timeout":   respTimeout,
		"send":               send,
		"ssl":                ssl,
	}
}

// NewPINGConfig constructs/returns a job configuration for PING type jobs.
// host is the IP address or hostname to ping. (Required)
// timeout is the timeout(in ms) before marking the host as failed.
// count is the number of packets to send.
// interval is the minimum time(in ms) to wait between sending each packet.
func NewPINGConfig(host string, timeout, count, interval int) *Config {
	return &Config{
		"host":     host, // Required
		"timeout":  timeout,
		"count":    count,
		"interval": interval,
	}
}
