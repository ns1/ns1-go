package mockns1

import (
	"net/http"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dataset"
)

// AddDatasetListTestCase sets up a test case for the api.Client.Datasets.List() function
func (s *Service) AddDatasetListTestCase(
	requestHeaders, responseHeaders http.Header,
	response []*dataset.Dataset,
) error {
	return s.AddTestCase(
		http.MethodGet, "/datasets", http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddDatasetGetTestCase sets up a test case for the api.Client.Datasets.Get() function
func (s *Service) AddDatasetGetTestCase(
	id string,
	requestHeaders, responseHeaders http.Header,
	response *dataset.Dataset,
) error {
	return s.AddTestCase(
		http.MethodGet, "/datasets/"+id, http.StatusOK, requestHeaders,
		responseHeaders, "", response,
	)
}

// AddDatasetCreateTestCase sets up a test case for the api.Client.Datasets.Create() function
func (s *Service) AddDatasetCreateTestCase(
	requestHeaders, responseHeaders http.Header,
	request, response *dataset.Dataset,
) error {
	return s.AddTestCase(
		http.MethodPut, "/datasets", http.StatusCreated, requestHeaders,
		responseHeaders, request, response,
	)
}

// AddDatasetDeleteTestCase sets up a test case for the api.Client.Datasets.Delete() function
func (s *Service) AddDatasetDeleteTestCase(
	id string, requestHeaders, responseHeaders http.Header,
) error {
	return s.AddTestCase(
		http.MethodDelete, "/datasets/"+id, http.StatusNoContent, requestHeaders,
		responseHeaders, "", "",
	)
}

// AddDatasetGetReportTestCase sets up a test case for the api.Client.Datasets.GetReport() function
func (s *Service) AddDatasetGetReportTestCase(
	id string, reportId string, requestHeaders, responseHeaders http.Header, fileContents []byte,
) error {
	return s.AddTestCase(
		http.MethodGet, "/datasets/"+id+"/reports/"+reportId, http.StatusOK, requestHeaders,
		responseHeaders, "", fileContents,
	)
}
