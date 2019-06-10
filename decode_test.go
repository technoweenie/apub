package apubencoding_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/technoweenie/apubencoding"
)

func TestDecode(t *testing.T) {
	t.Run("simple object", func(t *testing.T) {
		obj := Decode(t, `{
      "@context": "https://www.w3.org/ns/activitystreams",
      "type": "Object",
      "id": "http://www.test.example/object/1",
      "name": "A Simple, non-specific object"
    }`)

		assert.Equal(t, "https://www.w3.org/ns/activitystreams", obj.String("@context"))
		assert.Equal(t, obj.String("@context"), obj.Context())

		assert.Equal(t, "Object", obj.String("type"))
		assert.Equal(t, obj.String("type"), obj.Type())

		assert.Equal(t, "http://www.test.example/object/1", obj.String("id"))
		assert.Equal(t, obj.String("id"), obj.ID())

		assert.Equal(t, "A Simple, non-specific object", obj.String("name"))
		assert.Equal(t, obj.String("name"), obj.Name())
	})
}

func Decode(t *testing.T, input string) *apubencoding.Object {
	dec := &apubencoding.Decoder{}
	obj, err := dec.Decode(strings.NewReader(input))
	require.Nil(t, err)
	return obj
}
