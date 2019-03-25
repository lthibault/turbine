package turbine

import "sync/atomic"

type cursor struct {
	idx     int64
	padding [7]int64 // cache-line padding
}

func newCursor() *cursor { return &cursor{idx: -1} }

func (c *cursor) Load() int64 { return atomic.LoadInt64(&c.idx) }

func (c *cursor) Store(i int64) { atomic.StoreInt64(&c.idx, i) }

type barrier interface {
	Load() int64
}
