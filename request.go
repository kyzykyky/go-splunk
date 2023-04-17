package gosplunk

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Default insecure client
func getHttpClient() http.Client {
	return http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
}

// Do request and retry if needed
// Splunk sometimes return EOF error without reason
func (c Client) doRequest(request *http.Request) (*http.Response, error) {
	var response *http.Response
	var err error

	client := getHttpClient()
	client.Timeout = 10 * time.Second
	for i := 0; i < 3; i++ {
		response, err = client.Do(request)
		if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				c.Logger.Debug("request: EOF error, retrying...")
				continue
			}
			c.Logger.Warn("request:", err.Error())
			continue
		}
		break
	}
	return response, err
}

// Build common request for diffrent methods and users
func (c Client) requestBuilder(method string, useNs bool, resource string, query url.Values, body url.Values) (*http.Request, error) {
	if useNs {
		resource = c.getResourcePrefix() + resource
	}
	query.Add("output_mode", c.OutputMode)
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
	if c.Token != "" {
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}
	return request, nil
}

func responseReader(response *http.Response) ([]byte, error) {
	buf := &bytes.Buffer{}
	read, err := io.Copy(buf, response.Body)
	if err != nil {
		return []byte{}, err
	} else if read == 0 {
		return []byte{}, nil
	}
	return buf.Bytes(), nil
}

func (search NewSearch) setBody() url.Values {
	data := url.Values{}
	data.Set("search", search.Search)
	data.Set("earliest_time", search.Earliest)
	data.Set("latest_time", search.Latest)
	return data
}

func (jr SearchJobResultsRetrieve) setQuery() url.Values {
	params := url.Values{}
	params.Add("offset", strconv.Itoa(jr.Offset))
	params.Add("count", strconv.Itoa(jr.Count))
	return params
}

func (login Login) setBody() url.Values {
	data := url.Values{}
	data.Set("username", login.Username)
	data.Set("password", login.Password)
	return data
}

func (savedSearch SavedSearch) setBody() url.Values {
	data := url.Values{}
	if savedSearch.Name != "" {
		data.Set("name", savedSearch.Name)
	}
	if savedSearch.DisabledString != "" {
		data.Set("disabled", savedSearch.DisabledString)
	}
	if savedSearch.IsScheduledString != "" {
		data.Set("is_scheduled", savedSearch.IsScheduledString)
	}
	if savedSearch.Cron != "" {
		data.Set("cron_schedule", savedSearch.Cron)
	}
	if savedSearch.Description != "" {
		data.Set("description", savedSearch.Description)
	}
	if savedSearch.Search != "" {
		data.Set("search", savedSearch.Search)
	}
	if savedSearch.EarliestTime != "" {
		data.Set("dispatch.earliest_time", savedSearch.EarliestTime)
	}
	if savedSearch.LatestTime != "" {
		data.Set("dispatch.latest_time", savedSearch.LatestTime)
	}
	if savedSearch.RunOnStartupString != "" {
		data.Set("run_on_startup", savedSearch.RunOnStartupString)
	}
	if savedSearch.NextScheduledTime != "" {
		data.Set("next_scheduled_time", savedSearch.NextScheduledTime)
	}
	return data
}
