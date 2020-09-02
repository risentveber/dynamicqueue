package dynamicqueue

func NewDynamicallyBufferedChannel(in <- chan interface{}) <- chan interface{} {
	out := make(chan interface{})
	var storage []interface{}
	go func() {
		defer close(out)
		for {
			if len(storage) == 0 {
				item, ok := <-in
				if !ok {
					return
				}
				storage = append(storage, item)
				continue
			}

			select {
			case item, ok := <- in:
				if ok {
					storage = append(storage, item)
				} else {
					// unwind storage
					for _, item := range storage {
						out <- item
					}

					return
				}
			case out <- storage[0]:
				if len(storage) == 1 {
					storage = nil
				} else {
					storage = storage[1:]
				}
			}
		}
	}()
	return out
}
