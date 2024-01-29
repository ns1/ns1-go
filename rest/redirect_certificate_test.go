package rest_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/mockns1"
	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/redirect"
)

func TestRedirectCertificateService(t *testing.T) {
	mock, doer, err := mockns1.New(t)
	require.Nil(t, err)
	defer mock.Shutdown()

	client := api.NewClient(doer, api.SetEndpoint("https://"+mock.Address+"/v1/"))

	id := "87461586-c681-43b7-b283-f840be95e13c"

	t.Run("List", func(t *testing.T) {
		t.Run("Pagination", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = true
			certList := &redirect.CertificateList{
				Count: 4,
				Total: 4,
				Results: []*redirect.Certificate{
					{Domain: "a.com"},
					{Domain: "b.com"},
					{Domain: "c.com"},
					{Domain: "d.com"},
				},
			}
			require.Nil(t, mock.AddRedirectCertificateListTestCase(nil, nil, certList))

			respcerts, _, err := client.RedirectCertificates.List()
			require.Nil(t, err)
			require.NotNil(t, respcerts)
			require.Equal(t, len(certList.Results), len(respcerts))

			for i := range certList.Results {
				require.Equal(t, certList.Results[i].Domain, respcerts[i].Domain, i)
			}
		})

		t.Run("No Pagination", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = false
			certList := &redirect.CertificateList{
				Count: 2,
				Total: 4,
				Results: []*redirect.Certificate{
					{Domain: "a.com"},
					{Domain: "b.com"},
				},
			}

			header := http.Header{}
			header.Set("Link", `</redirect/certificates?after=b.com&limit=2>; rel="next"`)

			require.Nil(t, mock.AddRedirectCertificateListTestCase(nil, header, certList))

			respcerts, resp, err := client.RedirectCertificates.List()
			require.Nil(t, err)
			require.NotNil(t, respcerts)
			require.Equal(t, len(certList.Results), len(respcerts))
			require.Contains(t, resp.Header.Get("Link"), "redirect/certificates?after=b.com")

			for i := range certList.Results {
				require.Equal(t, certList.Results[i].Domain, respcerts[i].Domain, i)
			}
		})

		t.Run("Error", func(t *testing.T) {
			t.Run("HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "/redirect/certificates", http.StatusNotFound,
					nil, nil, "", `{"message": "test error"}`,
				))

				respcerts, resp, err := client.RedirectCertificates.List()
				require.Nil(t, respcerts)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "test error")
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			t.Run("Other", func(t *testing.T) {
				c := api.NewClient(errorClient{}, api.SetEndpoint(""))
				respcerts, resp, err := c.RedirectCertificates.List()
				require.Nil(t, resp)
				require.Error(t, err)
				require.Nil(t, respcerts)
			})
		})
	})

	t.Run("Get", func(t *testing.T) {

		t.Run("Pagination", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = true
			cert := &redirect.Certificate{
				Domain: "a.com",
			}
			require.Nil(t, mock.AddRedirectCertificateGetTestCase(id, nil, nil, cert))

			respcert, _, err := client.RedirectCertificates.Get(id)
			require.Nil(t, err)
			require.NotNil(t, respcert)
			require.Equal(t, cert.Domain, respcert.Domain)
		})

		t.Run("Error", func(t *testing.T) {
			t.Run("HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "/redirect/certificates/"+id, http.StatusNotFound,
					nil, nil, "", `{"message": "test error"}`,
				))

				cert, resp, err := client.RedirectCertificates.Get(id)
				require.Nil(t, cert)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "test error")
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			t.Run("Other", func(t *testing.T) {
				c := api.NewClient(errorClient{}, api.SetEndpoint(""))
				cert, resp, err := c.RedirectCertificates.Get(id)
				require.Nil(t, resp)
				require.Error(t, err)
				require.Nil(t, cert)
			})
		})
	})

	t.Run("Create", func(t *testing.T) {
		req := &redirect.Certificate{Domain: "a.com"}
		cert := &redirect.Certificate{
			ID:     &id,
			Domain: "a.com",
		}

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddRedirectCertificateCreateTestCase(nil, nil, req, cert))

			_, _, err := client.RedirectCertificates.Create(cert.Domain)
			require.Nil(t, err)
		})

		t.Run("Error", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodPut, "/redirect/certificates", http.StatusConflict,
				nil, nil, req, `{"message": "certificate already exists"}`,
			))

			_, _, err := client.RedirectCertificates.Create(cert.Domain)
			require.Equal(t, api.ErrRedirectCertificateExists, err)
		})
	})

	t.Run("Update", func(t *testing.T) {
		cert := &redirect.Certificate{
			ID:     &id,
			Domain: "a.com",
		}

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddRedirectCertificateUpdateTestCase(nil, nil, id, cert))

			_, err := client.RedirectCertificates.Update(id)
			require.Nil(t, err)
		})

		t.Run("Error", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodPost, "/redirect/certificates/"+id, http.StatusNotFound,
				nil, nil, "", `{"message": "certificate not found"}`,
			))

			_, err := client.RedirectCertificates.Update(id)
			require.Equal(t, api.ErrRedirectCertificateNotFound, err)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddRedirectCertificateDeleteTestCase(id, nil, nil))

			_, err := client.RedirectCertificates.Delete(id)
			require.Nil(t, err)
		})

		t.Run("Error", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodDelete, "/redirect/certificates/"+id, http.StatusNotFound,
				nil, nil, "", `{"message": "certificate not found"}`,
			))

			_, err := client.RedirectCertificates.Delete(id)
			require.Equal(t, api.ErrRedirectCertificateNotFound, err)
		})
	})
}
