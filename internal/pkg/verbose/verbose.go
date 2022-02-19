package verbose

import (
	"log"
	"sync"
)

var (
	isVerboseOn bool
	setOnce     sync.Once
)

func SetVerboseOn(verboseOn bool) {
	setOnce.Do(func() {
		isVerboseOn = verboseOn
	})
}

func Println(v ...interface{}) {
	if !isVerboseOn {
		return
	}

	log.Println(v...)
}

func Printf(format string, v ...interface{}) {
	if !isVerboseOn {
		return
	}

	log.Printf(format, v...)
}
