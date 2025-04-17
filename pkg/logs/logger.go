package logs

import (
	"errors"
	"log"
	"os"
)

var (
	Info           *log.Logger
	Error          *log.Logger
	SErrorNotFound error
	ErrorNotFound  string
)

func init() {
	Info = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	Error = log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
	SErrorNotFound = errors.New("Record not found")
	ErrorNotFound = "Record not found"
}
