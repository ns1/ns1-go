package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

func boolPtr(b bool) *bool {
	return &b
}

func TestCreateAPIKey(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		var k account.APIKey
		require.NoError(t, json.Unmarshal(b, &k))
		assert.Nil(t, k.Permissions.Security)
		assert.False(t, k.Permissions.Monitoring.ManageJobs)
		assert.False(t, k.Permissions.Monitoring.CreateJobs)
		assert.False(t, k.Permissions.Monitoring.UpdateJobs)
		assert.False(t, k.Permissions.Monitoring.DeleteJobs)

		_, err = w.Write(b)
		require.NoError(t, err)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	k := &account.APIKey{
		ID:          "id-1",
		Key:         "key-1",
		Name:        "name-1",
		Permissions: account.PermissionsMap{},
	}

	_, err := c.APIKeys.Create(k)
	require.NoError(t, err)
}

func TestCreateAPIKeyWithExpiryDuration(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/account/apikeys", r.URL.Path)

		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		var k account.APIKey
		require.NoError(t, json.Unmarshal(b, &k))
		assert.Equal(t, "30d", k.ExpiryDuration)
		assert.Equal(t, "test-key-with-expiry", k.Name)

		// Return the key with secrets populated
		response := account.APIKey{
			ID:             "id-123",
			Name:           k.Name,
			ExpiryDuration: k.ExpiryDuration,
			Permissions:    k.Permissions,
			Secrets: []*account.APIKeySecret{
				{
					ID:        "secret-abc123",
					Key:       "generated-secret-key",
					ExpiresAt: "2026-04-30T00:00:00Z",
					Enabled:   boolPtr(true),
				},
			},
		}

		respBytes, err := json.Marshal(response)
		require.NoError(t, err)
		_, err = w.Write(respBytes)
		require.NoError(t, err)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	k := &account.APIKey{
		Name:           "test-key-with-expiry",
		ExpiryDuration: "30d",
		Permissions:    account.PermissionsMap{},
	}

	_, err := c.APIKeys.Create(k)
	require.NoError(t, err)
	assert.Equal(t, "30d", k.ExpiryDuration)
	assert.NotNil(t, k.Secrets)
	assert.Len(t, k.Secrets, 1)
	assert.Equal(t, "secret-abc123", k.Secrets[0].ID)
	assert.Equal(t, "generated-secret-key", k.Secrets[0].Key)
}

func TestGetAPIKeyWithSecrets(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/account/apikeys/id-123", r.URL.Path)

		response := account.APIKey{
			ID:             "id-123",
			Name:           "test-key",
			ExpiryDuration: "30d",
			Permissions:    account.PermissionsMap{},
			Secrets: []*account.APIKeySecret{
				{
					ID:         "secret-1",
					ExpiresAt:  "2026-04-30T00:00:00Z",
					LastAccess: "2026-03-15T10:30:00Z",
					Enabled:    boolPtr(true),
				},
				{
					ID:        "secret-2",
					ExpiresAt: "2026-05-30T00:00:00Z",
					Enabled:   boolPtr(false),
				},
			},
		}

		respBytes, err := json.Marshal(response)
		require.NoError(t, err)
		_, err = w.Write(respBytes)
		require.NoError(t, err)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	k, _, err := c.APIKeys.Get("id-123")
	require.NoError(t, err)
	assert.Equal(t, "30d", k.ExpiryDuration)
	assert.NotNil(t, k.Secrets)
	assert.Len(t, k.Secrets, 2)
	assert.Equal(t, "secret-1", k.Secrets[0].ID)
	assert.True(t, *k.Secrets[0].Enabled)
	assert.Equal(t, "secret-2", k.Secrets[1].ID)
	assert.False(t, *k.Secrets[1].Enabled)
}

func TestUpdateSecret(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/apikeys/v1/secrets/secret-123", r.URL.Path)

		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		var secret account.APIKeySecret
		require.NoError(t, json.Unmarshal(b, &secret))
		assert.Equal(t, "secret-123", secret.ID)
		assert.False(t, *secret.Enabled)
		assert.Equal(t, "2026-06-01", secret.ExpiresAt)

		// Return updated secret
		response := account.APIKeySecret{
			ID:        "secret-123",
			ExpiresAt: "2026-06-01T00:00:00Z",
			Enabled:   boolPtr(false),
		}

		respBytes, err := json.Marshal(response)
		require.NoError(t, err)
		_, err = w.Write(respBytes)
		require.NoError(t, err)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	secret := &account.APIKeySecret{
		ID:        "secret-123",
		Enabled:   boolPtr(false),
		ExpiresAt: "2026-06-01",
	}

	_, err := c.APIKeys.UpdateSecret(secret)
	require.NoError(t, err)
	assert.Equal(t, "secret-123", secret.ID)
	assert.False(t, *secret.Enabled)
	assert.Equal(t, "2026-06-01T00:00:00Z", secret.ExpiresAt)
}

func TestUpdateSecretEnabledOnly(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)

		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		var secret account.APIKeySecret
		require.NoError(t, json.Unmarshal(b, &secret))
		assert.Equal(t, "secret-456", secret.ID)
		assert.True(t, *secret.Enabled)

		response := account.APIKeySecret{
			ID:        "secret-456",
			ExpiresAt: "2026-05-01T00:00:00Z",
			Enabled:   boolPtr(true),
		}

		respBytes, err := json.Marshal(response)
		require.NoError(t, err)
		_, err = w.Write(respBytes)
		require.NoError(t, err)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	secret := &account.APIKeySecret{
		ID:      "secret-456",
		Enabled: boolPtr(true),
	}

	_, err := c.APIKeys.UpdateSecret(secret)
	require.NoError(t, err)
	assert.True(t, *secret.Enabled)
}

func TestDeleteSecret(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/apikeys/v1/secrets/secret-789", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	_, err := c.APIKeys.DeleteSecret("secret-789")
	require.NoError(t, err)
}

func TestDeleteSecretMissing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		response := map[string]string{
			"message": "secret not found",
		}
		respBytes, _ := json.Marshal(response)
		_, _ = w.Write(respBytes)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	_, err := c.APIKeys.DeleteSecret("non-existent")
	require.Error(t, err)
	assert.Equal(t, ErrSecretMissing, err)
}

func TestUpdateSecretMissing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		response := map[string]string{
			"message": "secret not found",
		}
		respBytes, _ := json.Marshal(response)
		_, _ = w.Write(respBytes)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	secret := &account.APIKeySecret{
		ID:      "non-existent",
		Enabled: boolPtr(true),
	}

	_, err := c.APIKeys.UpdateSecret(secret)
	require.Error(t, err)
	assert.Equal(t, ErrSecretMissing, err)
}

func TestGetSecret(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/apikeys/v1/secrets/secret-123", r.URL.Path)

		response := account.APIKeySecret{
			ID:         "secret-123",
			ExpiresAt:  "2026-04-30T00:00:00Z",
			LastAccess: "2026-03-15T10:30:00Z",
			Enabled:    boolPtr(true),
		}

		respBytes, err := json.Marshal(response)
		require.NoError(t, err)
		_, err = w.Write(respBytes)
		require.NoError(t, err)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	secret, _, err := c.APIKeys.GetSecret("secret-123")
	require.NoError(t, err)
	assert.Equal(t, "secret-123", secret.ID)
	assert.True(t, *secret.Enabled)
	assert.Equal(t, "2026-04-30T00:00:00Z", secret.ExpiresAt)
	assert.Equal(t, "2026-03-15T10:30:00Z", secret.LastAccess)
}

func TestGetSecretMissing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		response := map[string]string{
			"message": "secret not found",
		}
		respBytes, _ := json.Marshal(response)
		_, _ = w.Write(respBytes)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	_, _, err := c.APIKeys.GetSecret("non-existent")
	require.Error(t, err)
	assert.Equal(t, ErrSecretMissing, err)
}

func TestGetSecretSelf(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/apikeys/v1/secrets/self", r.URL.Path)

		response := account.APIKeySecret{
			ID:         "secret-current",
			ExpiresAt:  "2026-05-15T00:00:00Z",
			LastAccess: "2026-03-20T14:22:00Z",
			Enabled:    boolPtr(true),
		}

		respBytes, err := json.Marshal(response)
		require.NoError(t, err)
		_, err = w.Write(respBytes)
		require.NoError(t, err)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	secret, _, err := c.APIKeys.GetSecretSelf()
	require.NoError(t, err)
	assert.Equal(t, "secret-current", secret.ID)
	assert.True(t, *secret.Enabled)
	assert.Equal(t, "2026-05-15T00:00:00Z", secret.ExpiresAt)
	assert.Equal(t, "2026-03-20T14:22:00Z", secret.LastAccess)
}

func TestGetSecretSelfInvalidAuth(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		response := map[string]string{
			"message": "invalid authentication credentials",
		}
		respBytes, _ := json.Marshal(response)
		_, _ = w.Write(respBytes)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	_, _, err := c.APIKeys.GetSecretSelf()
	require.Error(t, err)
	assert.Equal(t, ErrInvalidAuth, err)
}

func TestRenewSecret(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/apikeys/v1/secrets/secret-123/renew", r.URL.Path)

		// Renew returns the new secret with plaintext key
		response := account.APIKeySecret{
			ID:        "secret-new-456",
			Key:       "nss_AbCdEfGhIjKlMnOpQrStUvWxYz0987",
			ExpiresAt: "2026-06-13T00:00:00Z",
			Enabled:   boolPtr(true),
		}

		respBytes, err := json.Marshal(response)
		require.NoError(t, err)
		_, err = w.Write(respBytes)
		require.NoError(t, err)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	secret, _, err := c.APIKeys.RenewSecret("secret-123")
	require.NoError(t, err)
	assert.Equal(t, "secret-new-456", secret.ID)
	assert.Equal(t, "nss_AbCdEfGhIjKlMnOpQrStUvWxYz0987", secret.Key)
	assert.True(t, *secret.Enabled)
	assert.Equal(t, "2026-06-13T00:00:00Z", secret.ExpiresAt)
}

func TestRenewSecretMissing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		response := map[string]string{
			"message": "secret not found",
		}
		respBytes, _ := json.Marshal(response)
		_, _ = w.Write(respBytes)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	_, _, err := c.APIKeys.RenewSecret("non-existent")
	require.Error(t, err)
	assert.Equal(t, ErrSecretMissing, err)
}

func TestRenewSecretMaxSecretsReached(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]string{
			"message": "cannot renew secret: api key already has 2 active secrets",
		}
		respBytes, _ := json.Marshal(response)
		_, _ = w.Write(respBytes)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	_, _, err := c.APIKeys.RenewSecret("secret-123")
	require.Error(t, err)
	// Should return a generic error with the message
	assert.Contains(t, err.Error(), "cannot renew secret")
}

func TestRenewSecretSelf(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/apikeys/v1/secrets/self/renew", r.URL.Path)

		// Renew returns the new secret with plaintext key
		response := account.APIKeySecret{
			ID:        "secret-renewed-789",
			Key:       "nss_XyZ123AbC456DeF789GhI012JkL345",
			ExpiresAt: "2026-07-01T00:00:00Z",
			Enabled:   boolPtr(true),
		}

		respBytes, err := json.Marshal(response)
		require.NoError(t, err)
		_, err = w.Write(respBytes)
		require.NoError(t, err)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	secret, _, err := c.APIKeys.RenewSecretSelf()
	require.NoError(t, err)
	assert.Equal(t, "secret-renewed-789", secret.ID)
	assert.Equal(t, "nss_XyZ123AbC456DeF789GhI012JkL345", secret.Key)
	assert.True(t, *secret.Enabled)
	assert.Equal(t, "2026-07-01T00:00:00Z", secret.ExpiresAt)
}

func TestRenewSecretSelfInvalidAuth(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		response := map[string]string{
			"message": "invalid authentication credentials",
		}
		respBytes, _ := json.Marshal(response)
		_, _ = w.Write(respBytes)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	_, _, err := c.APIKeys.RenewSecretSelf()
	require.Error(t, err)
	assert.Equal(t, ErrInvalidAuth, err)
}

func TestRenewSecretSelfNoExpiryDuration(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]string{
			"message": "apikey expiry_duration is required for secret renewal",
		}
		respBytes, _ := json.Marshal(response)
		_, _ = w.Write(respBytes)
	}))
	defer ts.Close()
	c := NewClient(nil, SetEndpoint(ts.URL))

	_, _, err := c.APIKeys.RenewSecretSelf()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expiry_duration is required")
}
