package rest_test

import (
	"context"
	"net/http"
	"testing"

	"gopkg.in/ns1/ns1-go.v2/rest/model/pulsar"

	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/mockns1"
	api "gopkg.in/ns1/ns1-go.v2/rest"
)

func TestApplication(t *testing.T) {
	mock, doer, err := mockns1.New(t)
	require.Nil(t, err)
	defer mock.Shutdown()

	client := api.NewClient(doer, api.SetEndpoint("https://"+mock.Address+"/v1/"))

	t.Run("List", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			defer mock.ClearTestCases()
			client.FollowPagination = false

			applications := []*pulsar.Application{
				pulsar.NewApplication("app1"),
				pulsar.NewApplication("app2"),
			}

			header := http.Header{}
			header.Set("Link", `</applications?b.list.applications>`)

			require.Nil(t, mock.AddApplicationTestCase(nil, header, applications))

			respApplications, resp, err := client.Applications.List(context.Background())
			require.Nil(t, err)
			require.NotNil(t, respApplications)
			require.Equal(t, len(applications), len(respApplications))
			require.Contains(t, resp.Header.Get("Link"), "applications?b.list.applications")

			for i := range applications {
				require.Equal(t, applications[i].Name, respApplications[i].Name, i)
			}
		})

		t.Run("Error", func(t *testing.T) {
			t.Run("HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "/pulsar/apps", http.StatusNotFound,
					nil, nil, "", `{"message": "test error"}`,
				))

				applications, resp, err := client.Applications.List(context.Background())
				require.Nil(t, applications)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "test error")
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			t.Run("Other", func(t *testing.T) {
				c := api.NewClient(errorClient{}, api.SetEndpoint(""))
				applications, resp, err := c.Applications.List(context.Background())
				require.Nil(t, resp)
				require.Error(t, err)
				require.Nil(t, applications)
			})
		})
	})

	t.Run("Get", func(t *testing.T) {
		id := "a32fc"
		name := "MyApp"
		t.Run("success", func(t *testing.T) {
			defer mock.ClearTestCases()
			client.FollowPagination = false

			link := `pulsar/apps/` + id
			header := http.Header{}
			header.Set("Link", link)

			require.Nil(t, mock.AddApplicationGetTestCase(id,
				nil,
				header,
				pulsar.NewApplication(name)))

			respApplication, resp, err := client.Applications.Get(context.Background(), id)
			require.Nil(t, err)
			require.NotNil(t, respApplication)
			require.Equal(t, name, respApplication.Name)
			require.Contains(t, resp.Header.Get("Link"), link)
		})

		t.Run("Error", func(t *testing.T) {
			t.Run("HTTP", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, "/pulsar/apps/"+id, http.StatusNotFound,
					nil, nil, "", `{"message": "error text that will be ignored by API"}`,
				))

				respApplications, resp, err := client.Applications.Get(context.Background(), id)
				require.Nil(t, respApplications)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "does not exist")
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			t.Run("Other", func(t *testing.T) {
				c := api.NewClient(errorClient{}, api.SetEndpoint(""))
				applications, resp, err := c.Applications.Get(context.Background(), id)
				require.Nil(t, resp)
				require.Error(t, err)
				require.Nil(t, applications)
			})
		})
	})

	t.Run("Create", func(t *testing.T) {
		application := pulsar.NewApplication("App_test")

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddApplicationCreateTestCase(nil, nil, application, application))

			_, err := client.Applications.Create(context.Background(), application)
			require.Nil(t, err)
			require.Equal(t, application.Name, "App_test")
		})

		t.Run("Error", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodPut, "/pulsar/apps", http.StatusNotFound,
				nil, nil, application, `{"Message": "test error"}`,
			))

			_, err := client.Applications.Create(context.Background(), application)
			require.Contains(t, err.Error(), "test error")
		})
	})

	t.Run("Update", func(t *testing.T) {
		application := pulsar.NewApplication("App_test")
		application.ID = "a32fc"

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddApplicationUpdateTestCase(nil,
				nil,
				application,
				application))

			_, err := client.Applications.Update(context.Background(), application)
			require.Nil(t, err)
		})

		t.Run("Error", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodPost, "/pulsar/apps/"+application.ID, http.StatusNotFound,
				nil, nil, application, `{"message": "pulsar app not found"}`,
			))

			_, err := client.Applications.Update(context.Background(), application)
			require.Equal(t, api.ErrApplicationMissing, err)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		id := "a32fc1"
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddApplicationDeleteTestCase(id, nil, nil))

			_, err := client.Applications.Delete(context.Background(), id)
			require.Nil(t, err)
		})

		t.Run("Error", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddTestCase(
				http.MethodDelete, "/pulsar/apps/"+id, http.StatusNotFound,
				nil, nil, "", `{"message": "pulsar app not found"}`,
			))

			_, err := client.Applications.Delete(context.Background(), id)
			require.Equal(t, api.ErrApplicationMissing, err)
		})
	})
}
