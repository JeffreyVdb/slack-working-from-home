package util

import (
	"fmt"
	"io"
	"os"
)

func PrintError(w io.Writer, format string, args ...interface{}) {
	_, err := fmt.Fprintf(w, format, args...)
	if err != nil {
		panic(err.Error())
	}
}

func IsEnvDefined(envKey string) bool {
	_, ok := os.LookupEnv(envKey)
	return ok
}
