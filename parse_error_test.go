package apub_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/technoweenie/apub"
	"golang.org/x/xerrors"
)

func TestError(t *testing.T) {
	t.Run("valueAsObject", func(t *testing.T) {
		obj := Parse(t, `{
        "type": "Object",
        "object1": {
          "type": "TestObject",
          "badobject": 123,
          "test": "test"
        }
      }`)

		obj1, err := obj.FetchObject("object1")
		assert.Nil(t, err)

		missing, err := obj1.FetchObject("missing")
		assert.Nil(t, missing)
		assert.Nil(t, err)

		_, err = obj1.FetchObject("test")
		assert.True(t, xerrors.Is(err, apub.ErrKeyNotObject), err)

		_, err = obj1.FetchObject("badobject")
		assert.True(t, xerrors.Is(err, apub.ErrKeyTypeNotObject), err)

		i, err := obj1.FetchInt("test")
		assert.Equal(t, 0, i)
		if assert.NotNil(t, err) {
			assert.True(t, xerrors.Is(err, apub.ErrInvalidInt), err)
		}
	})

	t.Run("lang map", func(t *testing.T) {
		obj := Parse(t, `{
        "type": "Object",
        "name": "test",
        "image": {
          "type": "Image",
          "nameMap": {
            "en": "image"
          },
          "url": "http://example.com/image.jpg"
        }
      }`)

		name, err := obj.FetchLang("name", "en")
		assert.Equal(t, "test", name)
		assert.True(t, xerrors.Is(err, apub.ErrLangMapNotFound))

		img := obj.Object("image")
		require.NotNil(t, img)
		imgName, err := img.FetchLang("name", "en")
		assert.Equal(t, "image", imgName)
		assert.Nil(t, err)

		esName, err := img.FetchLang("name", "es")
		assert.Equal(t, "image", esName)
		assert.True(t, xerrors.Is(err, apub.ErrLangNotFound))
	})
}
