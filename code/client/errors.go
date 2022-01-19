package client

import "errors"

var errCount = errors.New("invalid count format")
var errSize = errors.New("invalid size format")
var errCR = errors.New("invalid cr")
var errLF = errors.New("invalid lf")
