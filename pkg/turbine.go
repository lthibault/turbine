package turbine

// Consumer is notified of indices in the ring buffer that can be modified safely
type Consumer interface {
	Consume(lower, upper int64)
}

// Turbine is an interface to a ring buffer
type Turbine struct {
	cap int64
	w   *Writer
	rs  []*reader
}

// Start the Turbine
func (t Turbine) Start() {
	for _, r := range t.rs {
		r.Start()
	}
}

// Stop the Turbine
func (t Turbine) Stop() {
	for _, r := range t.rs {
		r.Stop()
	}
}

// Writer for the turbine
func (t *Turbine) Writer() *Writer { return t.w }

// New Turbine
func New(cap int64, cs ...Consumer) (t *Turbine) {
	t = new(Turbine)
	t.cap = cap
	t.rs = make([]*reader, len(cs))

	cursors := make([]*cursor, len(cs)+1)
	for i := range cursors {
		cursors[i] = new(cursor)
	}

	var b barrier = cursors[0]
	written := cursors[0] // writer cursor

	for i, c := range cs {
		cr := cursors[i+1] // cursor 0 belongs to the writer
		t.rs[i] = newReader(cr, written, b, c)
		b = cr
	}

	t.w = newWriter(cap, written, b)

	return
}
