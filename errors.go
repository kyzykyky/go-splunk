package gosplunk

import (
	"errors"
)

type ErrGeneral struct {
	User string
}

func (e ErrGeneral) Error() string {
	return "unknown error"
}

var ErrEmptyCredentials error = errors.New("empty credentials")
var ErrAuthFailed error = errors.New("authentication failed")
var ErrRequest = errors.New("request error")
var ErrInvalidRequest error = errors.New("invalid request")
var ErrInvalidResponse error = errors.New("invalid response")
var ErrFailedAction error = errors.New("failed action")
var ErrForbiddenAction error = errors.New("forbidden action")
var ErrConflict error = errors.New("conflict")

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
	case 409:
		c.Logger.Warnw("conflict", "user", c.Username)
		return ErrConflict
	case 500:
		c.Logger.Errorw("action failed", "user", c.Username)
		return ErrFailedAction
	default:
		err := ErrGeneral{User: c.Username}
		c.Logger.Errorw(err.Error(), "user", c.Username, "status", status)
		return err
	}
}
