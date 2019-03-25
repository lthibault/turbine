package turbine

import (
	"testing"
)

const (
	cap  = 1024 // must be a power of 2
	mask = cap - 1
)

var (
	res  int
	seq  int64
	ring = [cap]packet{}
)

type packet struct {
	Value int
}

type global struct{}

func (global) Consume(lower, upper int64) {
	for seq := lower; seq <= upper; seq++ {
		res = ring[seq&mask].Value
	}
}

func BenchmarkTurbine(b *testing.B) {
	b.Run("ReserveOne", func(b *testing.B) {
		t := New(cap, global{})
		w := t.Writer()
		t.Start()

		for i := 0; i < b.N; i++ {
			seq = w.Reserve(1)
			ring[seq&mask].Value = i
			w.Commit(seq)
		}

		t.Stop()
	})

	const batch = 16
	b.Run("ReserveMany", func(b *testing.B) {
		t := New(cap, global{})
		w := t.Writer()
		t.Start()

		for i := 0; i < b.N; i++ {
			if i%batch == 0 {
				if seq != 0 {
					w.Commit(seq)
				}
				seq = w.Reserve(batch)
			}

			ring[seq&mask].Value = 1
		}

		t.Stop()
	})
}

func BenchmarkChan(b *testing.B) {
	b.Run("Unbuffered", func(b *testing.B) {
		ch := make(chan int)
		go func() {
			for i := range ch {
				res = i
			}
		}()

		for i := 0; i < b.N; i++ {
			ch <- i
		}

		close(ch)
	})

	b.Run("Buffered-1", func(b *testing.B) {
		ch := make(chan int, 1)
		go func() {
			for i := range ch {
				res = i
			}
		}()

		for i := 0; i < b.N; i++ {
			ch <- i
		}

		close(ch)
	})

	b.Run("Buffered-16", func(b *testing.B) {
		ch := make(chan int, 16)
		go func() {
			for i := range ch {
				res = i
			}
		}()

		for i := 0; i < b.N; i++ {
			ch <- i
		}

		close(ch)
	})

}
