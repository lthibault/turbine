package main

import (
	"context"
	"log"

	turbine "github.com/lthibault/turbine/pkg"
)

const (
	cap  = 1024 // must be a power of 2
	mask = cap - 1
)

var ring = [cap]packet{}

type packet struct {
	Value int
}

type logger struct{}

func (logger) Consume(lower, upper int64) {
	var msg packet
	for seq := lower; seq <= upper; seq++ {
		msg = ring[seq&mask]
		log.Println(msg.Value)
	}
}

func main() {
	t := turbine.New(cap, logger{})
	t.Start()
	defer t.Stop()

	go func() {
		w := t.Writer()

		for i := 0; i < 100; i++ {
			seq := w.Reserve(1)
			ring[seq&mask].Value = i
			w.Commit(seq)
		}
	}()

	<-context.Background().Done()
}
