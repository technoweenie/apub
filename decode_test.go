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
      "name": "A Simple, non-specific object",
			"url": "http://www.test.example/object/1.html"
    }`)

		assert.Equal(t, "https://www.w3.org/ns/activitystreams", obj.String("@context"))
		assert.Equal(t, obj.String("@context"), obj.Context())

		assert.Equal(t, "Object", obj.String("type"))
		assert.Equal(t, obj.String("type"), obj.Type())

		assert.Equal(t, "http://www.test.example/object/1", obj.String("id"))
		assert.Equal(t, obj.String("id"), obj.ID())

		assert.Equal(t, "A Simple, non-specific object", obj.String("name"))
		assert.Equal(t, obj.String("name"), obj.Name())

		assert.Equal(t, "http://www.test.example/object/1.html", obj.String("url"))

		link := obj.URL()
		require.NotNil(t, link)
		assert.Equal(t, "Link", link.String("type"))
		assert.Equal(t, link.String("type"), link.Type())
		assert.Equal(t, obj.String("url"), link.String("href"))
		assert.Equal(t, "", link.String("mediaType"))
	})

	t.Run("simple object with nested objects", func(t *testing.T) {
		obj := Decode(t, `{
      "@context": "https://www.w3.org/ns/activitystreams",
      "type": "Object",
      "id": "http://www.test.example/object/1",
			"url": {
				"type": "Link",
				"href": "http://www.test.example/object/1.html",
				"mediaType": "text/html"
			}
    }`)

		assert.Equal(t, "http://www.test.example/object/1.html", obj.String("url"))

		link := obj.URL()
		require.NotNil(t, link)
		assert.Equal(t, "Link", link.String("type"))
		assert.Equal(t, link.String("type"), link.Type())
		assert.Equal(t, obj.String("url"), link.String("href"))
		assert.Equal(t, "text/html", link.String("mediaType"))
	})
}

func Decode(t *testing.T, input string) *apubencoding.Object {
	dec := &apubencoding.Decoder{}
	obj, err := dec.Decode(strings.NewReader(input))
	require.Nil(t, err)
	return obj
}
