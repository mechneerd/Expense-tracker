package log

import (
	"log"
	"os"
)

var Std *log.Logger

func Init() {
	Std = log.New(os.Stdout, "expense-api: ", log.LstdFlags|log.Lshortfile)
}