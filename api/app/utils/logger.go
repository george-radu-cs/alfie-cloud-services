package utils

import (
	"log"
	"os"
)

var InfoLogger = log.New(
	os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile|log.LUTC|log.Lmsgprefix,
)
var ErrorLogger = log.New(
	os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile|log.LUTC|log.
		Lmsgprefix,
)
