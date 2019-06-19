package apub_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectSetKey(t *testing.T) {
	obj := Parse(t, `{}`)

	t.Run("bool", func(t *testing.T) {
		assert.False(t, obj.Bool("bool"))
		obj.SetBool("bool", true)
		assert.True(t, obj.Bool("bool"))
	})

	t.Run("num", func(t *testing.T) {
		assert.Equal(t, 0, obj.Int("int"))
		obj.SetNum("int", 108)
		assert.Equal(t, 108, obj.Int("int"))
	})

	t.Run("id list", func(t *testing.T) {
		assert.Equal(t, 0, len(obj.To()))
		obj.SetList("to", []interface{}{"http://example.com/1", "http://example.com/2"})
		to := obj.To()
		assert.Equal(t, 2, len(to))
		assert.Contains(t, to, "http://example.com/1")
		assert.Contains(t, to, "http://example.com/2")

		obj.AppendList("to", "http://example.com/3")
		to2 := obj.To()
		assert.Equal(t, 3, len(to2))
		assert.Contains(t, to2, "http://example.com/1")
		assert.Contains(t, to2, "http://example.com/2")
		assert.Contains(t, to2, "http://example.com/3")
	})

	t.Run("list", func(t *testing.T) {
		assert.Equal(t, 0, len(obj.IDs("attributedTo")))
		assert.Equal(t, 0, len(obj.List("attributedTo")))

		obj.SetList("attributedTo", []interface{}{
			"http://joe.example.org",
			map[string]interface{}{
				"type": "Person",
				"name": "Sally",
				"id":   "http://sally.example.org",
			},
		})

		t.Run("SetList + IDs()", func(t *testing.T) {
			ato := obj.AttributedTo()
			assert.Equal(t, 2, len(ato))
			assert.Contains(t, ato, "http://joe.example.org")
			assert.Contains(t, ato, "http://sally.example.org")
		})

		obj.AppendList("attributedTo",
			map[string]interface{}{
				"type": "Person",
				"name": "Bob",
				"id":   "http://bob.example.org",
			},
			"http://jane.example.org",
		)

		t.Run("AppendList + IDs()", func(t *testing.T) {
			ato := obj.AttributedTo()
			assert.Equal(t, 4, len(ato))
			assert.Contains(t, ato, "http://joe.example.org")
			assert.Contains(t, ato, "http://sally.example.org")
			assert.Contains(t, ato, "http://bob.example.org")
			assert.Contains(t, ato, "http://jane.example.org")
		})
	})

	t.Run("object", func(t *testing.T) {
		assert.Equal(t, "", obj.Object("target").Type())
		assert.Equal(t, "", obj.Str("target"))
		obj.SetObject("target", map[string]interface{}{
			"type": "Test",
			"id":   "http://example.com/test",
		})
		assert.Equal(t, "Test", obj.Object("target").Type())
		assert.Equal(t, "http://example.com/test", obj.Str("target"))
	})

	t.Run("string", func(t *testing.T) {
		assert.Equal(t, "", obj.Str("s"))
		obj.SetStr("s", "string")
		assert.Equal(t, "string", obj.Str("s"))
	})

	assert.Nil(t, obj.Errors())
	assert.Nil(t, obj.NonFatalErrors())
}

func TestObjectDeleteKey(t *testing.T) {
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
