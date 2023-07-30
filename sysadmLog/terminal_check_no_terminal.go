// +build js nacl plan9

package sysadmLog

import (
	"io"
)

func checkIfTerminal(w io.Writer) bool {
	return false
}
