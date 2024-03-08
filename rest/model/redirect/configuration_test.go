package redirect

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

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
	got := NewConfiguration(cfg.Domain, cfg.Path, cfg.Target, cfg.Tags, cfg.ForwardingMode, cfg.ForwardingType, cfg.HttpsEnabled, cfg.HttpsForced, cfg.QueryForwarding)
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
	https := true
	force := true
	queryFwd := true
	tag := "aaa"
	updated := time.Now().Unix()
	d := []byte(`{
		"id": "` + id + `",
		"certificate_id": "` + certId + `",
		"domain": "` + domain + `",
		"path": "` + path + `",
		"target": "` + target + `",
		"forwarding_mode": "` + fwMode.String() + `",
		"forwarding_type": "` + fwType.String() + `",
		"https_enabled": ` + strconv.FormatBool(https) + `,
		"https_forced": ` + strconv.FormatBool(force) + `,
		"query_forwarding": ` + strconv.FormatBool(queryFwd) + `,
		"tags": [
			"` + tag + `"
		],
		"last_updated": ` + strconv.FormatInt(updated, 10) + `
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
	assert.NotNil(t, cfg.HttpsEnabled, "nil https_enabled")
	assert.Equal(t, https, *cfg.HttpsEnabled, "https_enabled mismatch")
	assert.NotNil(t, cfg.HttpsForced, "nil https_forced")
	assert.Equal(t, force, *cfg.HttpsForced, "https_forced mismatch")
	assert.NotNil(t, cfg.QueryForwarding, "nil query_forwarding")
	assert.Equal(t, queryFwd, *cfg.QueryForwarding, "query_forwarding mismatch")
	assert.Len(t, cfg.Tags, 1, "tag size mismatch")
	assert.Equal(t, tag, cfg.Tags[0], "tags mismatch")
	assert.NotNil(t, cfg.LastUpdated, "nil last_updated")
	assert.Equal(t, updated, *cfg.LastUpdated, "last_updated mismatch")
}

func TestMarshalConfiguration(t *testing.T) {
	fwMode := None
	fwType := Masking
	https := true
	force := true
	queryFwd := true
	cfg := NewConfiguration(
		"www.mydomain.com",
		"/path",
		"https://google.com",
		[]string{"a", "b"},
		&fwMode,
		&fwType,
		&https,
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
		"https_enabled": `+strconv.FormatBool(https)+`,
		"https_forced": `+strconv.FormatBool(force)+`,
		"query_forwarding": `+strconv.FormatBool(queryFwd)+`,
		"tags": [
			"`+cfg.Tags[0]+`", "`+cfg.Tags[1]+`"
		]
	}`, string(d), "json mismatch")
}
