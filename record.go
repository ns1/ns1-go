package nsone

type Answer struct {
	Answer []string               `json:"answer,omitempty"`
	Meta   map[string]interface{} `json:"meta,omitempty"`
}

type Record struct {
	Id      string            `json:"id,omitempty"`
	Zone    string            `json:"zone,omitempty"`
	Domain  string            `json:"domain,omitempty"`
	Type    string            `json:"type,omitempty"`
	Link    string            `json:"link"`
	Meta    map[string]string `json:"meta,omitempty"`
	Answers []Answer          `json:"answers"`
}

func NewRecord(zone string, domain string, t string) *Record {
	return &Record{
		Zone:    zone,
		Domain:  domain,
		Type:    t,
		Answers: make([]Answer, 0),
	}
}

func (r *Record) LinkTo(to string) {
	r.Meta = nil
	r.Answers = make([]Answer, 0)
	r.Link = to
}
