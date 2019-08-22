package util

import (
	"fmt"
	"io"
)

func PrintError(w io.Writer, format string, args ...interface{}) {
	_, err := fmt.Fprintf(w, format, args...)
	if err != nil {
		panic(err.Error())
	}
}
