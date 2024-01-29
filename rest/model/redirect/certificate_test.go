package redirect

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCertificate(t *testing.T) {
	domain := "www.mydomain.com"
	cert := NewCertificate(domain)
	assert.Equal(t, &Certificate{Domain: domain}, cert, "certificate mismatch")
}

func TestUnmarshalCertificate(t *testing.T) {
	id := "313829b3-0861-44b1-b92a-cdd49d555a83"
	domain := "www.mydomain.com"
	certificate := "something"
	processing := false
	errors := "failed renewing certificate"
	d := []byte(`{
		"id": "` + id + `",
		"domain": "` + domain + `",
		"certificate": "` + certificate + `",
		"processing": ` + strconv.FormatBool(processing) + `,
		"errors": "` + errors + `"
	}`)
	cert := Certificate{}
	err := json.Unmarshal(d, &cert)
	assert.NoError(t, err, "unmarshalling error")
	assert.NotNil(t, cert.ID, "nil id")
	assert.Equal(t, id, *cert.ID, "id mismatch")
	assert.Equal(t, domain, cert.Domain, "domain mismatch")
	assert.NotNil(t, cert.Certificate, "nil certificate")
	assert.Equal(t, certificate, *cert.Certificate, "id mismatch")
	assert.NotNil(t, cert.Processing, "nil processing")
	assert.Equal(t, processing, *cert.Processing, "processing mismatch")
	assert.NotNil(t, cert.Errors, "nil errors")
	assert.Equal(t, errors, *cert.Errors, "errors mismatch")
}

func TestMarshalCertificate(t *testing.T) {
	domain := "www.mydomain.com"
	cert := NewCertificate(domain)
	d, err := json.Marshal(&cert)
	assert.NoError(t, err, "marshalling error")
	assert.JSONEq(t, `{
		"domain": "`+domain+`"
	}`, string(d), "json mismatch")
}
