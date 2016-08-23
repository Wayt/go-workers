package workers

import (
	"fmt"
	"runtime"
)

type errFatal struct {
	message string
	file    string
	line    int
}

func (e errFatal) Error() string {
	return fmt.Sprintf("%s:%d %s", e.file, e.line, e.message)
}

func Fatalf(message string, i ...interface{}) error {
	return Fatal(fmt.Sprintf(message, i...))
}

func Fatal(i ...interface{}) error {

	_, file, line, _ := runtime.Caller(1)
	return &errFatal{message: fmt.Sprint(i...), file: file, line: line}
}

func IsFatal(err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(*errFatal)
	return ok
}
