package guid

import (
	"fmt"
	"sync"
)

var mu sync.Mutex
var id int

func init() {
	id = -1
}

func New() string {
	mu.Lock()
	defer mu.Unlock()

	id++
	return fmt.Sprintf("s%v", id)
}
