package ns1

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
	path := fmt.Sprintf("data/sources")

	req, err := c.NewRequest("PUT", path, &ds)
	if err != nil {
		return err
	}

	// Update datasources fields with data from api(ensure consistent)
	_, err = c.Do(req, &ds)
	if err != nil {
		return err
	}

	return nil
}

// GetDataSource takes an ID returns the details for a single data source
func (c APIClient) GetDataSource(id string) (*DataSource, error) {
	path := fmt.Sprintf("data/sources/%s", id)
	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var ds DataSource
	_, err = c.Do(req, &ds)
	if err != nil {
		return nil, err
	}

	return &ds, nil
}

// DeleteDataSource takes an ID and removes an existing data source and all connected feeds from the cource
func (c APIClient) DeleteDataSource(id string) error {
	path := fmt.Sprintf("data/sources/%s", id)
	req, err := c.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

// UpdateDataSource takes a *DataSource modifies basic details of a data source
func (c APIClient) UpdateDataSource(ds *DataSource) error {
	path := fmt.Sprintf("data/sources/%s", ds.Id)

	req, err := c.NewRequest("POST", path, &ds)
	if err != nil {
		return err
	}

	// Update datasources fields with data from api(ensure consistent)
	_, err = c.Do(req, &ds)
	if err != nil {
		return err
	}

	return nil
}

// PublishFeed takes a datasources' id and data to publish to the feed.
func (c APIClient) PublishFeed(datasource_id string, data interface{}) error {
	path := fmt.Sprintf("feed/%s", datasource_id)

	req, err := c.NewRequest("POST", path, &data)
	if err != nil {
		return err
	}

	_, err = c.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}
