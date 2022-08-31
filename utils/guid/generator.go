package guid

import (
	"sync"
)

var mu sync.Mutex
var id int

func init() {
	id = 0
}

func New() int {
	mu.Lock()
	defer mu.Unlock()

	id++
	return id
}
