package error_list

import "errors"

var (
	ErrDuplicate = errors.New("duplicate")
)
