package util

import (
	"errors"
	"strings"
	"testing"
)

type ErrorProneNullWriter struct{}

func (w *ErrorProneNullWriter) Write(data []byte) (n int, err error) {
	if string(data) == "narwhals are not awesome" {
		return 0, errors.New("narwhals are awesome")
	}

	return
}

func TestPrintError(t *testing.T) {
	stringWriter := &strings.Builder{}
	PrintError(stringWriter, "some error")
	if stringWriter.String() != "some error" {
		t.Errorf("error output should be: some error")
	}
}

func TestPrintErrorFormat(t *testing.T) {
	stringWriter := &strings.Builder{}
	PrintError(stringWriter, "a should be less than %d", 5)
	if stringWriter.String() != "a should be less than 5" {
		t.Errorf("error should be equal to: \"a should be less than 5\"")
	}
}

func TestPrintErrorPanic(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("PrintError was supposed to panic")
		} else if err != "narwhals are awesome" {
			t.Errorf("Should have panicked with \"narwhals are awesome\"")
		}
	}()

	writer := &ErrorProneNullWriter{}
	PrintError(writer, "narwhals are not awesome")
}
