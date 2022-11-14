package gosplunk

type auth struct {
	SessionKey string          `json:"sessionKey"`
	Message    string          `json:"message"`
	Code       string          `json:"code"`
	Messages   []errorMessages `json:"messages,omitempty"`
}

type errorMessages struct {
	Type string `json:"type"`
	Code string `json:"code"`
	Text string `json:"text"`
}

type Login struct {
	Username string
	Password string
}

type Authorized struct {
	Username string
	Token    string
}

type User struct {
	Username string
	Realname string
	Roles    []string
}

type NewSearch struct {
	Search   string `json:"search"`
	Earliest string `json:"earliest"`
	Latest   string `json:"latest"`
}

type SearchJob struct {
	Job string `json:"sid" query:"sid"`
}

type JobInfo struct {
	Sid         string  `json:"sid"`
	IsDone      bool    `json:"isDone"`
	IsFailed    bool    `json:"isFailed"`
	EventCount  float64 `json:"eventCount"`
	ResultCount float64 `json:"resultCount"`
	RunDuration float64 `json:"runDuration"`
}
type JobResults struct {
	Results []map[string]interface{} `json:"results"`
	Offset  int                      `json:"offset"`
	Count   int                      `json:"count"`
}
type ExportJobResults struct {
	Result  map[string]interface{} `json:"result"`
	Lastrow bool                   `json:"lastrow"`
	Offset  int                    `json:"offset"`
}

type SearchJobResultsRetrieve struct {
	Job    string `json:"sid" query:"sid"`
	Count  int    `json:"count,omitempty"`
	Offset int    `json:"offset,omitempty"`
}

type SearchResults struct {
	Preview  bool                   `json:"preview"`
	Offset   int                    `json:"offset"`
	Lastrow  bool                   `json:"lastrow"`
	Result   map[string]interface{} `json:"result"`
	Messages []errorMessages        `json:"messages,omitempty"`
}
