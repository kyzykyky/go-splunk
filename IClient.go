package gosplunk

// Describes every action, that gosplunk can do
type IClient interface {
	// Search
	Search(search NewSearch) (SearchJob, error)
	SearchExport(search NewSearch) ([]ExportJobResults, error)
	JobAwait(job string, retries int) (JobInfo, error)
	JobResults(job SearchJobResultsRetrieve) (JobResults, error)
	JobRetrieve(job string) (JobInfo, error)

	// Lookups
	LookupRetrieve(search string) ([]map[string]interface{}, error)
	LookupRetrieveKv(search string) (map[string]string, error)

	// Saved searches
	SavedSearchGet(searchName string, ns NameSpace) (SavedSearch, error)
	SavedSearchCreate(search SavedSearch) error
	SavedSearchUpdate(search SavedSearch, ns NameSpace) error
	SavedSearchDelete(searchName string, ns NameSpace) error

	// Auth and check
	Ping() error
	SubClientByPass(user string, password string) Client
	SubClientByToken(user string, token string) Client
	UserAuth(login Login) (Authorized, error)
	UserInfo(username string) (User, error)
	UsersInfo() (map[string]User, error)
	ValidateToken(token string) error
}
