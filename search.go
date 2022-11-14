package splunk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (c Client) Search(search NewSearch) (SearchJob, error) {
	request, err := c.requestBuilder(http.MethodPost, "/search/jobs", url.Values{}, search.setBody())
	if err != nil {
		c.Logger.Error(err)
		return SearchJob{}, ErrRequest
	}
	response, err := getHttpClient().Do(request)
	if err != nil {
		c.Logger.Error(err)
		return SearchJob{}, ErrRequest
	} else if response.StatusCode >= 400 {
		return SearchJob{}, c.requestError(response.StatusCode)
	}
	defer response.Body.Close()

	var j SearchJob
	err = json.NewDecoder(response.Body).Decode(&j)
	if err != nil {
		c.Logger.Error(err)
		return SearchJob{}, ErrInvalidResponse
	}
	return j, nil
}

// Get job status and info
func (c Client) JobRetrieve(job string) (JobInfo, error) {
	request, err := c.requestBuilder(http.MethodGet, fmt.Sprintf("/search/jobs/%s", job), url.Values{}, url.Values{})
	if err != nil {
		c.Logger.Error(err)
		return JobInfo{}, ErrRequest
	}
	response, err := getHttpClient().Do(request)
	if err != nil {
		c.Logger.Error(err)
		return JobInfo{}, ErrRequest
	} else if response.StatusCode >= 400 {
		return JobInfo{}, c.requestError(response.StatusCode)
	}
	defer response.Body.Close()

	results := make(map[string]interface{})
	if err := json.NewDecoder(response.Body).Decode(&results); err != nil {
		c.Logger.Error(err)
		return JobInfo{}, ErrInvalidResponse
	}

	j, ok := results["entry"].([]interface{})[0].(map[string]interface{})["content"].(map[string]interface{})
	if !ok {
		c.Logger.Error(err, "results", results)
		return JobInfo{}, ErrInvalidResponse
	}
	jobInfo := JobInfo{Sid: job}

	if done, ok := j["isDone"].(bool); ok {
		jobInfo.IsDone = done
	}
	if failed, ok := j["isFailed"].(bool); ok {
		jobInfo.IsFailed = failed
	}
	if ecount, ok := j["eventCount"].(float64); ok {
		jobInfo.EventCount = ecount
	}
	if rcount, ok := j["resultCount"].(float64); ok {
		jobInfo.ResultCount = rcount
	}
	if runDuration, ok := j["runDuration"].(float64); ok {
		jobInfo.RunDuration = runDuration
	}
	return jobInfo, err
}

// Await job to finish
// Creates i request every 25 milliseconds awaiting for isDone=true
// TODO: redo with context
func (c Client) JobAwait(job string, retries int) (JobInfo, error) {
	var jobInfo JobInfo
	var err error
	for i := 0; i < retries; i++ {
		time.Sleep(25 * time.Millisecond)
		jobInfo, err = c.JobRetrieve(job)
		if err != nil {
			return jobInfo, err
		}
		if jobInfo.IsDone {
			break
		}
	}
	return jobInfo, nil
}

// Get results of job (max 100)
func (c Client) JobResults(job SearchJobResultsRetrieve) (JobResults, error) {
	request, err := c.requestBuilder(http.MethodGet, fmt.Sprintf("/search/jobs/%s/results", job.Job),
		job.setQuery(), url.Values{})
	if err != nil {
		c.Logger.Error(err)
		return JobResults{}, ErrRequest
	}
	response, err := getHttpClient().Do(request)

	if err != nil {
		c.Logger.Error(err)
		return JobResults{}, ErrRequest
	} else if response.StatusCode >= 400 {
		return JobResults{}, c.requestError(response.StatusCode)
	}
	defer response.Body.Close()
	results := make(map[string]interface{})
	err = json.NewDecoder(response.Body).Decode(&results)
	if err != nil {
		c.Logger.Error(err)
		return JobResults{}, ErrFailedAction
	}
	var count int
	if cs, ok := results["results"].([]interface{}); ok {
		count = len(cs)
	}

	res := make([]map[string]interface{}, count)
	if count != 0 {
		temp_res := results["results"].([]interface{})
		for i := range temp_res {
			res[i] = temp_res[i].(map[string]interface{})
		}
	}
	return JobResults{Results: res,
		Offset: job.Offset,
		Count:  count,
	}, nil
}

func (c Client) SearchExport(search NewSearch) ([]ExportJobResults, error) {
	request, err := c.requestBuilder(http.MethodPost, "/search/jobs/export",
		url.Values{}, search.setBody())
	if err != nil {
		c.Logger.Error(err)
		return []ExportJobResults{}, ErrRequest
	}
	response, err := getHttpClient().Do(request)
	if err != nil {
		c.Logger.Error(err)
		return []ExportJobResults{}, ErrRequest
	} else if response.StatusCode >= 400 {
		return []ExportJobResults{}, c.requestError(response.StatusCode)
	}
	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		c.Logger.Error(err)
	}
	bodyString := string(bodyBytes)
	searchResults := strings.Split(bodyString, "\n")

	exportResults := make([]ExportJobResults, 0)
	for _, item := range searchResults {
		var j ExportJobResults
		if len(item) > 0 {
			if err := json.Unmarshal([]byte(item), &j); err == nil {
				exportResults = append(exportResults, j)
			} else {
				c.Logger.Errorw(err.Error(), "item", item)
			}
		}
	}
	return exportResults, nil
}
