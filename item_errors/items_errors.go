package item_errors

import "errors"

var (
	RequestTimeoutErr = errors.New("request timeout")
	NotFoundErr = errors.New("item not found")
	ParseErr = errors.New("error when trying to parse response")
)
