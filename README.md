# Go-Splunk

Splunk-REST-API client for Go.

Initialize new client:

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
