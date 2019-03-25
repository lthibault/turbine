package turbine

import "runtime"

const (
	// spinMask represents the number of spins to wait before calling runtime.Gosched()
	// to avoid indefinite busy-loops.  This number is arbitrary.
	spinMask clutch = (1 << 14) - 1
)

type waitStrategy interface {
	Wait(int64)
}

type clutch int64

func (c clutch) Wait(spin int64) {
	if spin > 0 && clutch(spin)&c == 0 {
		runtime.Gosched()
	}
}

// Writer tracks the addition of elements to the ring buffer
type Writer struct {
	written         *cursor
	b               barrier
	cap, prev, gate int64
	s               waitStrategy
}

func newWriter(cap int64, c *cursor, b barrier) *Writer {
	return &Writer{
		cap:     cap,
		written: c,
		b:       b,
		s:       spinMask,
	}
}

// Reserve slots on the ring buffer
func (w *Writer) Reserve(n int64) int64 {
	w.prev += n

	for spin := int64(0); w.prev-w.cap > w.gate; spin++ {
		w.s.Wait(spin)
		w.gate = w.b.Load()
	}

	return w.prev
}

// Commit slots on the ring buffer
func (w *Writer) Commit(n int64) { w.written.Store(n) }
