package turbine

import (
	"runtime"
	"time"
)

type reader struct {
	read, written *cursor
	b             barrier
	c             Consumer
	ready         bool
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
	r.ready = true
	r.startReceiving()
}

func (r *reader) Stop() { r.ready = false }

func (r *reader) startReceiving() {
	go func() {
		var lower, upper, prev int64
		prev = r.read.Load()

		for {
			lower = prev + 1
			upper = r.b.Load()

			if lower <= upper {
				r.c.Consume(lower, upper)
				r.read.Store(upper)
				prev = upper
			} else if upper = r.written.Load(); lower <= upper {
				// N.B.:  sleeping increases the batch size by allowing the ring buffer
				//		  to fill.  This in turn reduces the number of writes required
				//		  to store the sequence. Reducing writes allows the CPU to
				//		  optimize the pipeline by avoiding prediction failures.
				time.Sleep(time.Microsecond)
			} else if r.ready {
				time.Sleep(time.Millisecond)
			} else {
				break
			}

			runtime.Gosched()
		}
	}()
}
