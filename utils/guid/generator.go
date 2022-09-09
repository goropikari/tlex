package guid

import (
	"sync"
)

var mu sync.Mutex
var id int

func init() {
	id = -1
}

func New() int {
	mu.Lock()
	defer mu.Unlock()

	id++
	return id
}
