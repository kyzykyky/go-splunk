package gosplunk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Basic saved search fields
type SavedSearch struct {
	DisabledString     string `json:"-"`
	Disabled           bool   `json:"disabled,omitempty"`
	Name               string `json:"name,omitempty"`
	IsScheduledString  string `json:"-"`
	IsScheduled        bool   `json:"is_scheduled,omitempty"`
	Cron               string `json:"cron_schedule,omitempty"`
	Description        string `json:"description,omitempty"`
	Search             string `json:"search,omitempty"`
	EarliestTime       string `json:"dispatch.earliest_time,omitempty"`
	LatestTime         string `json:"dispatch.latest_time,omitempty"`
	RunOnStartupString string `json:"-"`
	RunOnStartup       bool   `json:"run_on_startup,omitempty"`
	NextScheduledTime  string `json:"next_scheduled_time,omitempty"`
}

func (c Client) SavedSearchCreate(search SavedSearch, ns NameSpace) error {
	resource := fmt.Sprintf("%s/saved/searches", getResourcePrefix(ns))
	request, err := c.requestBuilder(http.MethodPost, false, resource, url.Values{}, search.setBody())
	if err != nil {
		c.Logger.Error(err)
		return ErrRequest
	}

	response, err := getHttpClient().Do(request)
	if err != nil {
		c.Logger.Error(err)
		return ErrRequest
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		resp, err := responseReader(response)
		if err == nil {
			c.Logger.Debug(string(resp))
		}
		return c.requestError(response.StatusCode)
	}
	return nil
}

// Update saved search with not empty struct fields
func (c Client) SavedSearchUpdate(search SavedSearch, ns NameSpace) error {
	resource := fmt.Sprintf("%s/saved/searches/%s", getResourcePrefix(ns), search.Name)
	search.Name = "" // Set name to empty string to avoid error
	request, err := c.requestBuilder(http.MethodPost, false, resource, url.Values{}, search.setBody())
	if err != nil {
		c.Logger.Error(err)
		return ErrRequest
	}
	response, err := getHttpClient().Do(request)
	if err != nil {
		c.Logger.Error(err)
		return ErrRequest
	}
	defer response.Body.Close()
	if response.StatusCode >= 400 {
		resp, err := responseReader(response)
		if err == nil {
			c.Logger.Debug(string(resp))
		}
		return c.requestError(response.StatusCode)
	}
	return nil
}

// Get saved search by name
func (c Client) SavedSearchGet(searchName string, ns NameSpace) (SavedSearch, error) {
	resource := fmt.Sprintf("%s/saved/searches/%s", getResourcePrefix(ns), searchName)
	request, err := c.requestBuilder(http.MethodGet, false, resource, url.Values{}, url.Values{})
	if err != nil {
		c.Logger.Error(err)
		return SavedSearch{}, ErrRequest
	}
	response, err := getHttpClient().Do(request)
	if err != nil {
		c.Logger.Error(err)
		return SavedSearch{}, ErrRequest
	}
	defer response.Body.Close()

	resp, err := responseReader(response)
	if response.StatusCode >= 400 {
		if err == nil {
			c.Logger.Debug(string(resp))
		}
		return SavedSearch{}, c.requestError(response.StatusCode)
	}

	var rawResponse map[string]interface{}
	err = json.Unmarshal(resp, &rawResponse)
	if err != nil {
		c.Logger.Errorw("saved search pre unmarshall error", "error", err.Error())
		return SavedSearch{}, ErrInvalidResponse
	}

	content := rawResponse["entry"].([]interface{})[0].(map[string]interface{})["content"].(map[string]interface{})
	contentbytes, err := json.Marshal(content)
	if err != nil {
		c.Logger.Errorw("saved search post unmarshall error", "error", err.Error())
		return SavedSearch{}, ErrInvalidResponse
	}

	var savedSearch SavedSearch
	err = json.Unmarshal(contentbytes, &savedSearch)
	if err != nil {
		c.Logger.Errorw("saved search unmarshall error", "error", err.Error())
		return SavedSearch{}, ErrInvalidResponse
	}
	return savedSearch, nil
}

// Delete saved search by name
func (c Client) SavedSearchDelete(searchName string, ns NameSpace) error {
	resource := fmt.Sprintf("%s/saved/searches/%s", getResourcePrefix(ns), searchName)
	request, err := c.requestBuilder(http.MethodDelete, false, resource, url.Values{}, url.Values{})
	if err != nil {
		c.Logger.Error(err)
		return ErrRequest
	}
	response, err := getHttpClient().Do(request)
	if err != nil {
		c.Logger.Error(err)
		return ErrRequest
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		stresp, err := responseReader(response)
		if err == nil {
			c.Logger.Debug(stresp)
		}
		return c.requestError(response.StatusCode)
	}
	return nil
}
