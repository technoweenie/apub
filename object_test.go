package apub_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	obj := Parse(t, `{
		"a": 1,
		"b": 2,
		"inner": {
			"c": 3,
			"d": 4
		}
	}`)

	t.Run("property", func(t *testing.T) {
		assert.Equal(t, 1, obj.Int("a"))
		assert.Equal(t, 2, obj.Int("b"))
		obj.Del("b")
		assert.Equal(t, 1, obj.Int("a"))
		assert.Equal(t, 0, obj.Int("b"))
	})

	t.Run("inner object", func(t *testing.T) {
		inner := obj.Object("inner")
		assert.Equal(t, 3, inner.Int("c"))
		assert.Equal(t, 4, inner.Int("d"))
		inner.Del("c")
		assert.Equal(t, 0, inner.Int("c"))
		assert.Equal(t, 4, inner.Int("d"))

		inner2 := obj.Object("inner")
		assert.Equal(t, 0, inner2.Int("c"))
		assert.Equal(t, 4, inner2.Int("d"))

		obj.Del("inner")
		inner3 := obj.Object("inner")
		assert.Equal(t, 0, inner3.Int("c"))
		assert.Equal(t, 0, inner3.Int("d"))
	})
}
