package turbine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCursor(t *testing.T) {
	c := new(cursor)

	t.Run("Store", func(t *testing.T) {
		if c.Store(9001); !assert.Equal(t, int64(9001), c.idx) {
			t.FailNow()
		}
	})

	t.Run("Load", func(t *testing.T) {
		if !assert.Equal(t, int64(9001), c.Load()) {
			t.FailNow()
		}
	})
}
