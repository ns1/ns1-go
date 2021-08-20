package rest_test

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ns1/ns1-go.v2/mockns1"

	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/pulsar"
)

func TestPulsarJob(t *testing.T) {
	mock, doer, err := mockns1.New(t)

	require.Nil(t, err)
	defer mock.Shutdown()

	client := api.NewClient(doer, api.SetEndpoint("https://"+mock.Address+"/v1/"))

	// Tests for api.Client.PulsarJobs.List()
	t.Run("List", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			pulsarJobs := []*pulsar.PulsarJob{
				{
					Name:   "PulsarJob1",
					TypeID: "latency",
					Config: pulsar.JobConfig{
						Host:     "myHost1",
						URL_Path: "/myURLPath1",
					},
				},
				{
					Name:   "PulsarJob2",
					TypeID: "custom",
				},
			}

			require.Nil(t, mock.AddPulsarJobListTestCase(myAppID, nil, nil, pulsarJobs))

			respPulsarJobs, _, err := client.PulsarJobs.List(myAppID)
			require.Nil(t, err)
			require.NotNil(t, respPulsarJobs)
			require.Equal(t, len(pulsarJobs), len(respPulsarJobs))

			for i := range pulsarJobs {
				require.Equal(t, pulsarJobs[i].Name, respPulsarJobs[i].Name, i)
				require.Equal(t, pulsarJobs[i].TypeID, respPulsarJobs[i].TypeID, i)
				if respPulsarJobs[i].TypeID == "latency" {
					require.Equal(t, pulsarJobs[i].Config.Host, respPulsarJobs[i].Config.Host, i)
					require.Equal(t, pulsarJobs[i].Config.URL_Path, respPulsarJobs[i].Config.URL_Path, i)
				}
			}
		})

		t.Run("Error", func(t *testing.T) {
			// Error Application not found
			t.Run("Application not found", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, fmt.Sprintf("/pulsar/apps/%v/jobs", myAppID), http.StatusNotFound,
					nil, nil, "", `{"message": "pulsar app not found"}`,
				))

				pulsarJobs, resp, err := client.PulsarJobs.List(myAppID)
				require.Nil(t, pulsarJobs)
				require.NotNil(t, err)
				require.Equal(t, api.ErrAppMissing.Error(), err.Error())
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			// Other errors
			t.Run("Other errors", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, fmt.Sprintf("/pulsar/apps/%s/jobs", myAppID), http.StatusNotFound,
					nil, nil, "", `{"message": "test error"}`,
				))

				pulsarJobs, resp, err := client.PulsarJobs.List(myAppID)
				require.Nil(t, pulsarJobs)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "test error")
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})
		})
	})

	// Tests for api.Client.PulsarJobs.Get()
	t.Run("Get", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			pulsarJob := &pulsar.PulsarJob{
				Customer:  2156,
				TypeID:    "latency",
				Name:      "myPulsarJob",
				Community: true,
				JobID:     myJobID,
				AppID:     myAppID,
				Active:    false,
				Shared:    true,
				Config: pulsar.JobConfig{
					Host:                 "myHost",
					URL_Path:             "/myURLPath",
					Http:                 false,
					Https:                true,
					RequestTimeoutMillis: 1234,
					JobTimeoutMillis:     4321,
					UseXHR:               true,
					StaticValues:         true,
					BlendMetricWeights: &pulsar.BlendMetricWeights{
						Timestamp: 567,
						Weights: []*pulsar.Weights{
							{
								Name:         "myWeight2",
								Weight:       12.3,
								DefaultValue: 32.1,
								Maximize:     true,
							},
							{
								Name:         "myWeight1",
								Weight:       12.3,
								DefaultValue: 32.1,
								Maximize:     false,
							},
						},
					},
				},
			}

			require.Nil(t, mock.AddPulsarJobGetTestCase(myAppID, myJobID, nil, nil, pulsarJob))

			respPulsarJob, _, err := client.PulsarJobs.Get(myAppID, myJobID)

			require.Nil(t, err)
			require.True(t, reflect.DeepEqual(pulsarJob, respPulsarJob))
		})

		t.Run("Error", func(t *testing.T) {
			// Error Job not found
			t.Run("Job not found", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, fmt.Sprintf("/pulsar/apps/%s/jobs/%s", myAppID, myJobID), http.StatusNotFound,
					nil, nil, "", fmt.Sprintf(`{"message": "pulsar job %s not found for appid %s"}`, myJobID, myAppID),
				))
				pulsarJob, resp, err := client.PulsarJobs.Get(myAppID, myJobID)
				require.Nil(t, pulsarJob)
				require.NotNil(t, err)
				require.Equal(t, api.ErrJobMissing.Error(), err.Error())
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			// Error Application not found
			t.Run("Application not found", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, fmt.Sprintf("/pulsar/apps/%s/jobs/%s", myAppID, myJobID), http.StatusNotFound,
					nil, nil, "", `{"message": "pulsar app not found"}`,
				))
				pulsarJob, resp, err := client.PulsarJobs.Get(myAppID, myJobID)
				require.Nil(t, pulsarJob)
				require.NotNil(t, err)
				require.Equal(t, api.ErrAppMissing.Error(), err.Error())
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			// Other errors
			t.Run("Other errors", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodGet, fmt.Sprintf("/pulsar/apps/%s/jobs/%s", myAppID, myJobID), http.StatusNotFound,
					nil, nil, "", `{"message": "test error"}`,
				))

				pulsarJob, resp, err := client.PulsarJobs.Get(myAppID, myJobID)
				require.Nil(t, pulsarJob)
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "test error")
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})
		})
	})

	t.Run("Create", func(t *testing.T) {
		pulsarJob := &pulsar.PulsarJob{
			Name:   "myPulsarJob",
			TypeID: "latency",
			AppID:  myAppID,
		}

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddPulsarJobCreateTestCase(nil, nil, pulsarJob, pulsarJob))

			_, err := client.PulsarJobs.Create(pulsarJob)
			require.Nil(t, err)
		})

		t.Run("Error", func(t *testing.T) {
			// Application not found
			t.Run("Application not found", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodPut, fmt.Sprintf("pulsar/apps/%s/jobs", myAppID), http.StatusNotFound,
					nil, nil, pulsarJob, `{"message": "pulsar app not found"}`,
				))

				_, err = client.PulsarJobs.Create(pulsarJob)

				require.Equal(t, api.ErrAppMissing.Error(), err.Error())
			})

			// Other errors
			t.Run("Other errors", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodPut, fmt.Sprintf("pulsar/apps/%s/jobs", myAppID), http.StatusNotFound,
					nil, nil, pulsarJob, `{"message": "test error"}`,
				))

				_, err = client.PulsarJobs.Create(pulsarJob)

				require.Contains(t, err.Error(), "test error")
			})
		})
	})

	t.Run("Update", func(t *testing.T) {
		pulsarJob := &pulsar.PulsarJob{
			Name:   "updatedPulsarJob",
			TypeID: "custom",
			AppID:  myAppID,
			JobID:  myJobID,
		}

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddPulsarJobUpdateTestCase(nil, nil, pulsarJob, pulsarJob))

			_, err := client.PulsarJobs.Update(pulsarJob)
			require.Nil(t, err)
		})

		t.Run("Error", func(t *testing.T) {
			// Error Job not found
			t.Run("Job not found", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodPost, fmt.Sprintf("/pulsar/apps/%s/jobs/%s", pulsarJob.AppID, pulsarJob.JobID), http.StatusNotFound,
					nil, nil, pulsarJob, fmt.Sprintf(`{"message": "pulsar job %s not found for appid %s"}`, myJobID, myAppID),
				))
				resp, err := client.PulsarJobs.Update(pulsarJob)
				require.NotNil(t, err)
				require.Equal(t, api.ErrJobMissing.Error(), err.Error())
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			// Error Application not found
			t.Run("Application not found", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodPost, fmt.Sprintf("/pulsar/apps/%s/jobs/%s", pulsarJob.AppID, pulsarJob.JobID), http.StatusNotFound,
					nil, nil, pulsarJob, `{"message": "pulsar app not found"}`,
				))
				resp, err := client.PulsarJobs.Update(pulsarJob)

				require.NotNil(t, err)
				require.Equal(t, api.ErrAppMissing.Error(), err.Error())
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})

			// Other errors
			t.Run("Other errors", func(t *testing.T) {
				defer mock.ClearTestCases()

				require.Nil(t, mock.AddTestCase(
					http.MethodPost, fmt.Sprintf("/pulsar/apps/%s/jobs/%s", pulsarJob.AppID, pulsarJob.JobID), http.StatusNotFound,
					nil, nil, pulsarJob, `{"message": "test error"}`,
				))
				resp, err := client.PulsarJobs.Update(pulsarJob)

				require.NotNil(t, err)
				require.Contains(t, err.Error(), "test error")
				require.Equal(t, http.StatusNotFound, resp.StatusCode)
			})
		})
	})

	t.Run("Delete", func(t *testing.T) {
		pulsarJob := &pulsar.PulsarJob{
			Name:   "myPulsarJob",
			TypeID: "custom",
			AppID:  myAppID,
			JobID:  myJobID,
		}

		t.Run("Success", func(t *testing.T) {
			defer mock.ClearTestCases()

			require.Nil(t, mock.AddPulsarJobDeleteTestCase(nil, nil, pulsarJob, nil))

			_, err := client.PulsarJobs.Delete(pulsarJob)
			require.Nil(t, err)
		})

		t.Run("Error", func(t *testing.T) {
			t.Run("Error", func(t *testing.T) {
				// Error Job not found
				t.Run("Job not found", func(t *testing.T) {
					defer mock.ClearTestCases()

					require.Nil(t, mock.AddTestCase(
						http.MethodDelete, fmt.Sprintf("/pulsar/apps/%s/jobs/%s", pulsarJob.AppID, pulsarJob.JobID), http.StatusNotFound,
						nil, nil, "", fmt.Sprintf(`{"message": "pulsar job %s not found for appid %s"}`, pulsarJob.JobID, pulsarJob.AppID),
					))
					resp, err := client.PulsarJobs.Delete(pulsarJob)
					require.NotNil(t, err)
					require.Equal(t, api.ErrJobMissing.Error(), err.Error())
					require.Equal(t, http.StatusNotFound, resp.StatusCode)
				})

				// Error Application not found
				t.Run("Application not found", func(t *testing.T) {
					defer mock.ClearTestCases()

					require.Nil(t, mock.AddTestCase(
						http.MethodDelete, fmt.Sprintf("/pulsar/apps/%s/jobs/%s", pulsarJob.AppID, pulsarJob.JobID), http.StatusNotFound,
						nil, nil, "", `{"message": "pulsar app not found"}`,
					))
					resp, err := client.PulsarJobs.Delete(pulsarJob)

					require.NotNil(t, err)
					require.Equal(t, api.ErrAppMissing.Error(), err.Error())
					require.Equal(t, http.StatusNotFound, resp.StatusCode)
				})

				// Other errors
				t.Run("Other errors", func(t *testing.T) {
					defer mock.ClearTestCases()

					require.Nil(t, mock.AddTestCase(
						http.MethodDelete, fmt.Sprintf("/pulsar/apps/%s/jobs/%s", pulsarJob.AppID, pulsarJob.JobID), http.StatusNotFound,
						nil, nil, "", `{"message": "test error"}`,
					))
					resp, err := client.PulsarJobs.Delete(pulsarJob)

					require.NotNil(t, err)
					require.Contains(t, err.Error(), "test error")
					require.Equal(t, http.StatusNotFound, resp.StatusCode)
				})
			})

		})
	})
}

var (
	myAppID = "myAppID"
	myJobID = "myJobID"
)
