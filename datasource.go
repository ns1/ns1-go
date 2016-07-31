package nsone

import "fmt"

// FeedDestination wraps an element of a DataSource's "destinations" attribute
type FeedDestination struct {
	Destid   string `json:"destid"`
	Desttype string `json:"desttype"`
	Record   string `json:"record"`
}

// DataSource wraps an NS1 /data/sources resource
type DataSource struct {
	Id           string            `json:"id,omitempty"`
	Name         string            `json:"name"`
	SourceType   string            `json:"sourcetype"`
	Config       map[string]string `json:"config,omitempty"`
	Status       string            `json:"status,omitempty"`
	Destinations []FeedDestination `json:"destinations,omitempty"`
}

// NewDataSource takes a name and sourceType and creates a new *DataSource
func NewDataSource(name string, sourceType string) *DataSource {
	cf := make(map[string]string, 0)
	return &DataSource{
		Name:       name,
		SourceType: sourceType,
		Config:     cf,
	}
}

// CreateDataSource takes a *DataSource and creates a new data source
func (c APIClient) CreateDataSource(ds *DataSource) error {
	return c.doHTTPBoth("PUT", "https://api.nsone.net/v1/data/sources", ds)
}

// GetDataSource takes an ID returns the details for a single data source
func (c APIClient) GetDataSource(id string) (*DataSource, error) {
	ds := DataSource{}
	_, err := c.doHTTPUnmarshal("GET", fmt.Sprintf("https://api.nsone.net/v1/data/sources/%s", id), nil, &ds)
	return &ds, err
}

// DeleteDataSource takes an ID and removes an existing data source and all connected feeds from the cource
func (c APIClient) DeleteDataSource(id string) error {
	return c.doHTTPDelete(fmt.Sprintf("https://api.nsone.net/v1/data/sources/%s", id))
}

// UpdateDataSource takes a *DataSource modifies basic details of a data source
func (c APIClient) UpdateDataSource(ds *DataSource) error {
	return c.doHTTPBoth("POST", fmt.Sprintf("https://api.nsone.net/v1/data/sources/%s", ds.Id), ds)
}
