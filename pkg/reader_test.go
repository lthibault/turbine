package turbine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlag(t *testing.T) {
	f := new(flag)

	t.Run("FalseOnInit", func(t *testing.T) {
		if !assert.False(t, f.Bool()) {
			t.FailNow()
		}
	})

	t.Run("Set", func(t *testing.T) {
		if f.Set(); !assert.True(t, f.Bool()) {
			t.FailNow()
		}
	})

	t.Run("Unset", func(t *testing.T) {
		if f.Unset(); !assert.False(t, f.Bool()) {
			t.FailNow()
		}
	})
}
