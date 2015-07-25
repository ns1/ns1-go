package nsone

type DataSource struct {
	Id string `json:"id,omitempty"`
}

func NewDataSource() *DataSource {
	return &DataSource{}
}
