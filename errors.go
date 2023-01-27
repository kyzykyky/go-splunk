package gosplunk

import (
	"errors"
)

// TODO: Specify error types as structs implementing the error interface

var ErrGeneral = errors.New("general error")
var ErrEmptyCredentials = errors.New("empty credentials")
var ErrAuthFailed = errors.New("authentication failed")
var ErrRequest = errors.New("request error")
var ErrInvalidRequest = errors.New("invalid request")
var ErrInvalidResponse = errors.New("invalid response")
var ErrFailedAction = errors.New("failed action")
var ErrForbiddenAction = errors.New("forbidden action")
var ErrConflict = errors.New("conflict")
var ErrNotFound = errors.New("object not found")
var ErrParseRegex = errors.New("regex does not match")

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
		return ErrNotFound
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
		c.Logger.Errorw("unknown error", "user", c.Username, "status", status)
		return ErrGeneral
	}
}
