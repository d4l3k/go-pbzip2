package pbzip2

import (
	"log"
	"sync"
)

var warnOnce sync.Once

func warn() {
	warnOnce.Do(func() {
		log.Println("WARNING: missing pbzip2 from system")
	})
}
