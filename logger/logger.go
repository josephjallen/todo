package logger

import (
	"log"
	"os"
)

type TraceIdKey struct{}

var InfoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
var WarningLog = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)

// When logging error messages it is good practice to use 'os.Stderr' instead of os.Stdout
var ErrorLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
