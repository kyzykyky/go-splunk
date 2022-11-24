# Go-Splunk

Splunk-REST-API client for Go.

## Initialize new client

```go
gosplunk.NewClient(gosplunk.ClientConfig{
    Host:          "https://127.0.0.1:8089",
    Username:      "admin",
    Token:         "*******************",
    // Optional to use password instead of token, then a session token will be retrieved automatically
    Password:      "changeme",

    // Optional Splunk app context
    App:           "search",
    EnableLogging: false,
})
```

## Currently available client actions

```go
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
```

Examples are available in _test.go files.
