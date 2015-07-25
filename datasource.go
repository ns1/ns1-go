package nsone

type DataSource struct {
	Id         string `json:"id,omitempty"`
	Name       string `json:"name"`
	SourceType string `json"source_type"`
}

func NewDataSource(name string, source_type string) *DataSource {
	return &DataSource{
		Name:       name,
		SourceType: source_type,
	}
}
