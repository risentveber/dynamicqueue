package dynamicqueue

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

type handler struct {
	sync.Mutex
	names []string
}

func (h *handler) OnItem(item interface{}) {
	time.Sleep(time.Second)
	h.Lock()
	defer h.Unlock()
	fmt.Println("process", item)
	h.names = append(h.names, item.(string))
}

func TestNewDynamicQueue(t *testing.T) {
	h := &handler{}
	q := NewDynamicQueue(3, 3, h)
	for i := 0; i < 10; i++ {
		q.Add(strconv.Itoa(i))
		fmt.Println("add", i)
	}

	q.Stop()
	if len(h.names) != 10 {
		t.Errorf("invalid count %d %s", len(h.names), h.names)
	}
}
