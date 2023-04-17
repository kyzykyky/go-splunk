package gosplunk

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
)

type ClientConfig struct {
	Host          string
	App           string
	SubApp        string
	Username      string
	Password      string
	Token         string
	EnableLogging bool
	Logger        Logger
	OutputMode    string
}

type Client struct {
	ClientConfig
	Roles []string
}

func NewClient(c ClientConfig) (Client, error) {
	C := Client{
		ClientConfig: c,
	}
	if c.Host == "" {
		C.Host = os.Getenv("SPLUNK_HOST")
	}
	if c.Username == "" {
		C.Username = os.Getenv("SPLUNK_ADMIN_USERNAME")
	}
	if c.Password == "" {
		C.Password = os.Getenv("SPLUNK_ADMIN_PASSWORD")
	}
	if c.App == "" {
		C.App = os.Getenv("SPLUNK_APP")
	}
	if c.Token == "" {
		C.Token = os.Getenv("SPLUNK_TOKEN")
		if C.Token == "" {
			st, err := Auth(&C)
			if err != nil || st == "" {
				return C, err
			}
			C.Token = st
		}
	}
	if !c.EnableLogging {
		c.Logger = NoLogger{}
	}
	if c.OutputMode == "" {
		// Default output mode
		c.OutputMode = "json"
	}
	if err := C.Ping(); err != nil {
		return C, err
	}
	return C, nil
}

// Create request as another client
func (c *Client) SubClientByPass(user, password string) Client {
	auth, err := c.UserAuth(Login{user, password})
	if err != nil {
		c.Logger.Errorw(err.Error(), "user", user)
		return Client{}
	}
	return Client{
		ClientConfig: ClientConfig{
			Host:          c.Host,
			App:           c.App,
			Username:      auth.Username,
			Password:      password,
			Token:         auth.Token,
			EnableLogging: c.EnableLogging,
			Logger:        c.Logger,
		},
	}
}

func (c *Client) SubClientByToken(user, token string) Client {
	return Client{
		ClientConfig: ClientConfig{
			Host:          c.Host,
			App:           c.App,
			Username:      user,
			Token:         token,
			EnableLogging: c.EnableLogging,
			Logger:        c.Logger,
		},
	}
}

func Auth(c *Client) (string, error) {
	login := Login{
		Username: c.Username,
		Password: c.Password,
	}
	if login.Username == "" || login.Password == "" {
		return "", ErrEmptyCredentials
	}

	request, err := c.requestBuilder(http.MethodPost, false, "/services/auth/login", url.Values{}, login.setBody())
	if err != nil {
		c.Logger.Error(err)
		return "", ErrRequest
	}

	response, err := c.doRequest(request)
	// TODO: Better error handling
	if err != nil {
		c.Logger.Error(err, "user", login.Username)
		return "", ErrRequest
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		c.Logger.Infow(ErrAuthFailed.Error(), "status", response.StatusCode, "user", login.Username)
		stresp, err := responseReader(response)
		if err == nil {
			c.Logger.Debug(stresp)
		}
		return "", ErrAuthFailed
	}

	authresp := auth{}
	err = json.NewDecoder(response.Body).Decode(&authresp)
	if err != nil {
		c.Logger.Error(err)
		return "", ErrInvalidResponse
	}
	if authresp.SessionKey != "" {
		return authresp.SessionKey, nil
	} else {
		c.Logger.Warnw(ErrAuthFailed.Error(), "response", authresp.Messages[0].Text)
		return "", ErrAuthFailed
	}
}

func (c Client) Ping() error {
	request, err := c.requestBuilder(http.MethodGet, false, "/services/server/info", url.Values{}, url.Values{})
	if err != nil {
		c.Logger.Error(err)
		return ErrRequest
	}
	response, err := c.doRequest(request)
	if err != nil {
		c.Logger.Error(err)
		return ErrRequest
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		c.Logger.Info(ErrAuthFailed)
		return ErrAuthFailed
	}
	return nil
}
