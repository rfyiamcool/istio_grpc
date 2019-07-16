package constant

import (
	"errors"
	"os"
)

const (
	OK    = "ok"
	Error = "error"
)

var (
	HostName string

	ErrNotFound = errors.New("key not found")
)

func init() {
	HostName, _ = os.Hostname()
}
