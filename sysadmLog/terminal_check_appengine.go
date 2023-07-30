// +build appengine

package sysadmLog

import (
	"io"
)

func checkIfTerminal(w io.Writer) bool {
	return true
}
