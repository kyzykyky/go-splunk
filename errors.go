package gosplunk

import (
	"errors"
)

var ErrGeneral = errors.New("unknown error")
var ErrEmptyCredentials error = errors.New("empty credentials")
var ErrAuthFailed error = errors.New("authentication failed")
var ErrRequest = errors.New("request error")
var ErrInvalidRequest error = errors.New("invalid request")
var ErrInvalidResponse error = errors.New("invalid response")
var ErrFailedAction error = errors.New("failed action")
var ErrForbiddenAction error = errors.New("forbidden action")

var ErrJobNotFound error = errors.New("job not found")

func (c Client) requestError(status int) error {
	switch status {
	case 200:
		return nil
	case 400:
		c.Logger.Errorw("bad request", "user", c.Username)
		return ErrInvalidRequest
	case 401:
		c.Logger.Errorw("auth failed", "user", c.Username)
		return ErrAuthFailed
	case 403:
		c.Logger.Warnw("action forbidden", "user", c.Username)
		return ErrForbiddenAction
	case 404:
		return ErrJobNotFound
	case 405:
		c.Logger.Errorw("method not allowed", "user", c.Username)
		return ErrInvalidRequest
	case 500:
		c.Logger.Errorw("action failed", "user", c.Username)
		return ErrFailedAction
	default:
		c.Logger.Errorw("unknown error", "user", c.Username)
		return ErrGeneral
	}
}
