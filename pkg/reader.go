package turbine

import (
	"runtime"
	"sync/atomic"
	"time"
	"unsafe"
)

type flag uint32

func (f *flag) Bool() bool { return atomic.LoadUint32((*uint32)(unsafe.Pointer(f))) != 0 }

func (f *flag) Set() {
	atomic.CompareAndSwapUint32((*uint32)(unsafe.Pointer(f)), 0, 1)
}
func (f *flag) Unset() {
	atomic.CompareAndSwapUint32((*uint32)(unsafe.Pointer(f)), 1, 0)
}

type reader struct {
	read, written *cursor
	b             barrier
	c             Consumer
	ready         flag
}

func newReader(read, written *cursor, b barrier, c Consumer) *reader {
	return &reader{
		read:    read,
		written: written,
		b:       b,
		c:       c,
	}
}

func (r *reader) Start() {
	r.ready.Set()
	go r.receive()
}

func (r *reader) Stop() { r.ready.Unset() }

func (r *reader) receive() {
	previous := r.read.Load()
	var lower, upper int64

	for {
		lower = previous + 1
		upper = r.b.Load()

		if lower <= upper {
			r.c.Consume(lower, upper)
			r.read.Store(upper)
			previous = upper
		} else if upper = r.written.Load(); lower <= upper {
			time.Sleep(time.Microsecond)
		} else if r.ready.Bool() {
			time.Sleep(time.Millisecond)
		} else {
			break
		}

		// sleeping increases the batch size which reduces number of writes required to
		// store the sequence reducing the number of writes allows the CPU to optimize
		// the pipeline without prediction failures
		runtime.Gosched()
	}
}
