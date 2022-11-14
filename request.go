package splunk

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Default insecure client
var httpClient *http.Client = &http.Client{Transport: &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}}

func getHttpClient() *http.Client {
	return httpClient
}

// Build common request for diffrent methods and users
func (c Client) requestBuilder(method string, resource string, query url.Values, body url.Values) (*http.Request, error) {
	resource = c.getResourcePrefix() + resource

	// Add default output_mode to query params
	query.Add("output_mode", "json")

	endcodedBody := body.Encode()

	u, err := url.ParseRequestURI(c.Host)
	if err != nil {
		return nil, err
	}
	u.Path = resource
	u.RawQuery = query.Encode()
	urlStr := u.String()

	request, err := http.NewRequest(method, urlStr, strings.NewReader(endcodedBody))
	if err != nil {
		return nil, err
	}

	// Set headers
	if method == http.MethodPost {
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		request.Header.Add("Content-Length", strconv.Itoa(len(endcodedBody)))
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	return request, nil
}

// Struct - url.Values converter
// type urlValues interface {
// 	queryValues
// 	bodyValues
// }
// type queryValues interface {
// 	setQuery() url.Values
// }
// type bodyValues interface {
// 	setBody() url.Values
// }

func (search NewSearch) setBody() url.Values {
	data := url.Values{}
	data.Set("search", search.Search)
	data.Set("earliest", search.Earliest)
	data.Set("latest", search.Latest)
	return data
}

func (jr SearchJobResultsRetrieve) setQuery() url.Values {
	params := url.Values{}
	params.Add("offset", strconv.Itoa(jr.Offset))
	params.Add("count", strconv.Itoa(jr.Count))
	return params
}
