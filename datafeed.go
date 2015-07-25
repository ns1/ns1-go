package nsone

/*
{
	  "destinations": [ ],
	    "id": "5208fd909f782d3d253d62ae",
		  "data": { },
		    "config": {
				    "entity_id": "enAAAAA",
					    "alarm_id": "a1AAAAA"
						  },
						    "name": "server.example.com alert"
						}
*/
type DataFeed struct {
	SourceId string            `json:"-"`
	Id       string            `json:"id,omitempty"`
	Name     string            `json:"name"`
	Config   map[string]string `json:"config,omitempty"`
	Data     map[string]string `json:"data,omitempty"`
}

func NewDataFeed(source_id string) *DataFeed {
	return &DataFeed{SourceId: source_id}
}
