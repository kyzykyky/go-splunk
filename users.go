package gosplunk

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (c *Client) UserAuth(login Login) (Authorized, error) {
	resource := "/services/auth/login"

	params := url.Values{}
	params.Add("output_mode", "json")

	data := url.Values{}
	data.Set("username", login.Username)
	data.Set("password", login.Password)

	u, _ := url.ParseRequestURI(c.Host)
	u.Path = resource
	u.RawQuery = params.Encode()
	urlStr := u.String()

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	request, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode()))

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	response, err := client.Do(request)
	if err != nil {
		c.Logger.Error(err, "username", login.Username)
		return Authorized{}, ErrRequest
	} else if response.StatusCode >= 400 {
		c.Logger.Info(ErrAuthFailed, "username", login.Username)
		return Authorized{}, ErrAuthFailed
	}
	defer response.Body.Close()
	authresp := auth{}
	err = json.NewDecoder(response.Body).Decode(&authresp)
	if err != nil {
		c.Logger.Error(err, "username", login.Username)
		return Authorized{}, ErrInvalidResponse
	}
	if authresp.SessionKey != "" {
		return Authorized{login.Username, authresp.SessionKey}, nil
	} else {
		c.Logger.Warnw(ErrAuthFailed.Error(), "username", login.Username, "response", authresp.Messages[0].Text)
		return Authorized{}, ErrAuthFailed
	}
}

// Check token
func (c *Client) ValidateToken(token string) error {
	resource := "/services/server/info"

	params := url.Values{}
	data := url.Values{}

	u, _ := url.ParseRequestURI(c.Host)
	u.Path = resource
	u.RawQuery = params.Encode()
	urlStr := u.String()

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	request, _ := http.NewRequest(http.MethodGet, urlStr, strings.NewReader(data.Encode()))

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	response, err := client.Do(request)
	if err != nil {
		c.Logger.Error(err)
		return ErrRequest
	} else if response.StatusCode >= 400 {
		c.Logger.Info(ErrAuthFailed)
		return ErrAuthFailed
	}
	return nil
}

// Get user info
// TODO: Redo with direct API call
func (c *Client) UserInfo(username string) (User, error) {

	resource := c.getResourcePrefix() + "/search/jobs/export"

	params := url.Values{}
	params.Add("output_mode", "json")

	data := url.Values{}
	data.Set("search", fmt.Sprintf("| rest /services/authentication/users splunk_server=local | table title realname roles | where title=\"%s\"", username))

	u, _ := url.ParseRequestURI(c.Host)
	u.Path = resource
	u.RawQuery = params.Encode()
	urlStr := u.String()

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	request, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode()))

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	response, err := client.Do(request)
	if err != nil {
		c.Logger.Error(err)
		return User{}, ErrRequest
	}

	defer response.Body.Close()
	sr := SearchResults{}
	if err := json.NewDecoder(response.Body).Decode(&sr); err != nil {
		c.Logger.Error(err)
		return User{}, ErrInvalidResponse
	}

	// Auth or user not found errors
	// TODO: Split handling
	if response.StatusCode >= 400 {
		err := errors.New(sr.Messages[0].Text)
		c.Logger.Warn(ErrInvalidResponse.Error(), "response", err)
		return User{}, ErrFailedAction
	}

	// Convert roles to []string
	roles := multivalueConverter(sr.Result["roles"])

	realname, ok := sr.Result["realname"].(string)
	if !ok {
		c.Logger.Warn("realname is not a string")
	}

	return User{username, realname, roles}, nil
}

func (c *Client) UsersInfo() (map[string]User, error) {
	res, err := c.RetrieveLookup("| rest /services/authentication/users splunk_server=local | table title realname roles")
	if err != nil {
		return nil, err
	}
	users := make(map[string]User)
	for _, user := range res {
		roles := multivalueConverter(user["roles"])
		users[user["title"].(string)] = User{
			Username: user["title"].(string),
			Realname: user["realname"].(string),
			Roles:    roles,
		}
	}
	return users, nil
}
