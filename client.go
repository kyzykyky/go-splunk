package gosplunk

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/kyzykyky/go-splunk/logger"
)

type ClientConfig struct {
	Host          string
	App           string
	Username      string
	Password      string
	Token         string
	EnableLogging bool
	Logger        logger.Logger
}

type Client struct {
	ClientConfig
	Roles []string
}

func NewClient(c ClientConfig) Client {
	C := Client{
		ClientConfig: c,
	}
	if c.Host == "" {
		C.Host = os.Getenv("SPLUNK_HOST")
	}
	if c.Username == "" {
		C.Username = os.Getenv("SPLUNK_USER")
	}
	if c.Password == "" {
		C.Password = os.Getenv("SPLUNK_PASSWORD")
	}
	if c.App == "" {
		C.App = os.Getenv("SPLUNK_APP")
	}
	if c.Token == "" {
		C.Token = os.Getenv("SPLUNK_TOKEN")
		if C.Token == "" {
			err := Auth(&C)
			if err != nil {
				panic(err)
			}
		}
	}
	if !c.EnableLogging {
		c.Logger = logger.NoLogger{}
	}
	if err := C.Ping(); err != nil {
		panic(err)
	}
	return C
}

func Auth(c *Client) error {
	if c.Username == "" || c.Password == "" {
		return ErrEmptyCredentials
	}
	resource := "/services/auth/login"

	params := url.Values{}
	params.Add("output_mode", "json")

	data := url.Values{}
	data.Set("username", c.Username)
	data.Set("password", c.Password)

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
	// TODO: Better error handling
	if err != nil {
		c.Logger.Error(err)
		return ErrRequest
	} else if response.StatusCode >= 400 {
		c.Logger.Info(ErrAuthFailed)
		return ErrAuthFailed
	}
	defer response.Body.Close()

	authresp := auth{}
	err = json.NewDecoder(response.Body).Decode(&authresp)
	if err != nil {
		c.Logger.Error(err)
		return ErrInvalidResponse
	}
	if authresp.SessionKey != "" {
		c.Token = authresp.SessionKey
		return nil
	} else {
		c.Logger.Warnw(ErrAuthFailed.Error(), "response", authresp.Messages[0].Text)
		return ErrAuthFailed
	}
}

func (c Client) Ping() error {
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

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
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
