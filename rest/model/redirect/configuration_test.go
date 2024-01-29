package redirect

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfiguration(t *testing.T) {
	fwMode := All
	fwType := Permanent
	cfg := Configuration{
		Domain:         "www.mydomain.com",
		Path:           "/path",
		Target:         "https://google.com",
		ForwardingMode: &fwMode,
		ForwardingType: &fwType,
	}
	got := NewConfiguration(cfg.Domain, cfg.Path, cfg.Target, cfg.Tags, cfg.ForwardingMode, cfg.ForwardingType, cfg.SslEnabled, cfg.ForceRedirect, cfg.QueryForwarding)
	assert.Equal(t, &cfg, got, "Configuration mismatch")
}

func TestNewConfigurationMinimal(t *testing.T) {
	cfg := Configuration{
		Domain: "www.mydomain.com",
		Path:   "/path",
		Target: "https://google.com",
	}
	got := NewConfigurationMinimal(cfg.Domain, cfg.Path, cfg.Target)
	assert.Equal(t, &cfg, got, "Configuration mismatch")
}

func TestUnmarshalConfiguration(t *testing.T) {
	id := "313829b3-0861-44b1-b92a-cdd49d555a83"
	certId := "7b3c6d8b-721b-4a61-9520-84cfae07cb94"
	domain := "www.mydomain.com"
	path := "/*"
	target := "https://www.google.com"
	fwMode := Capture
	fwType := Temporary
	ssl := true
	force := true
	queryFwd := true
	tag := "aaa"
	d := []byte(`{
		"id": "` + id + `",
		"certificate_id": "` + certId + `",
		"domain": "` + domain + `",
		"path": "` + path + `",
		"target": "` + target + `",
		"forwarding_mode": "` + fwMode.String() + `",
		"forwarding_type": "` + fwType.String() + `",
		"ssl_enabled": ` + strconv.FormatBool(ssl) + `,
		"force_redirect": ` + strconv.FormatBool(force) + `,
		"query_forwarding": ` + strconv.FormatBool(queryFwd) + `,
		"tags": [
			"` + tag + `"
		]
	}`)
	cfg := Configuration{}
	err := json.Unmarshal(d, &cfg)
	assert.NoError(t, err, "unmarshalling error")
	assert.NotNil(t, cfg.ID, "nil id")
	assert.Equal(t, id, *cfg.ID, "id mismatch")
	assert.NotNil(t, cfg.CertificateID, "nil certificate id")
	assert.Equal(t, certId, *cfg.CertificateID, "certificate_id mismatch")
	assert.Equal(t, domain, cfg.Domain, "domain mismatch")
	assert.Equal(t, path, cfg.Path, "path mismatch")
	assert.Equal(t, target, cfg.Target, "target mismatch")
	assert.NotNil(t, cfg.ForwardingMode, "nil forwarding_mode")
	assert.Equal(t, fwMode, *cfg.ForwardingMode, "forwarding_mode mismatch")
	assert.NotNil(t, cfg.ForwardingType, "nil forwarding_type")
	assert.Equal(t, fwType, *cfg.ForwardingType, "forwarding_type mismatch")
	assert.NotNil(t, cfg.SslEnabled, "nil ssl_enabled")
	assert.Equal(t, ssl, *cfg.SslEnabled, "ssl_enabled mismatch")
	assert.NotNil(t, cfg.ForceRedirect, "nil force_redirect")
	assert.Equal(t, force, *cfg.ForceRedirect, "force_redirect mismatch")
	assert.NotNil(t, cfg.QueryForwarding, "nil query_forwarding")
	assert.Equal(t, queryFwd, *cfg.QueryForwarding, "query_forwarding mismatch")
	assert.Len(t, cfg.Tags, 1, "tag size mismatch")
	assert.Equal(t, tag, cfg.Tags[0], "tags mismatch")
}

func TestMarshalConfiguration(t *testing.T) {
	fwMode := None
	fwType := Masking
	ssl := true
	force := true
	queryFwd := true
	cfg := NewConfiguration(
		"www.mydomain.com",
		"/path",
		"https://google.com",
		[]string{"a", "b"},
		&fwMode,
		&fwType,
		&ssl,
		&force,
		&queryFwd,
	)
	d, err := json.Marshal(&cfg)
	assert.NoError(t, err, "marshalling error")
	assert.JSONEq(t, `{
		"domain": "`+cfg.Domain+`",
		"path": "`+cfg.Path+`",
		"target": "`+cfg.Target+`",
		"forwarding_mode": "`+cfg.ForwardingMode.String()+`",
		"forwarding_type": "`+cfg.ForwardingType.String()+`",
		"ssl_enabled": `+strconv.FormatBool(ssl)+`,
		"force_redirect": `+strconv.FormatBool(force)+`,
		"query_forwarding": `+strconv.FormatBool(queryFwd)+`,
		"tags": [
			"`+cfg.Tags[0]+`", "`+cfg.Tags[1]+`"
		]
	}`, string(d), "json mismatch")
}
