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

	t.Run("object url", func(t *testing.T) {
		assertLink := func(t *testing.T, link *apubencoding.Object, expHref, expMediaType string) {
			require.NotNil(t, link)
			assert.Equal(t, "Link", link.String("type"))
			assert.Equal(t, "Link", link.Type())
			assert.Equal(t, expHref, link.String("href"))
			assert.Equal(t, expMediaType, link.String("mediaType"))
		}

		pdfURL := "http://example.org/4q-sales-forecast.pdf"
		pdfType := "application/pdf"

		t.Run("string", func(t *testing.T) {
			obj := Decode(t, `{
				"@context": "https://www.w3.org/ns/activitystreams",
				"type": "Object",
				"name": "4Q Sales Forecast",
				"url": "http://example.org/4q-sales-forecast.pdf"
			}`)

			links := obj.URLs()
			require.Equal(t, 1, len(links))
			assertLink(t, links[0], pdfURL, "")
			assert.Equal(t, pdfURL, obj.String("url"))
		})

		t.Run("subclass string", func(t *testing.T) {
			obj := Decode(t, `{
				"@context": "https://www.w3.org/ns/activitystreams",
				"type": "Document",
				"name": "4Q Sales Forecast",
				"url": "http://example.org/4q-sales-forecast.pdf"
			}`)

			links := obj.URLs()
			require.Equal(t, 1, len(links))
			assertLink(t, links[0], pdfURL, "")
			assert.Equal(t, pdfURL, obj.String("url"))
		})

		t.Run("subclass object", func(t *testing.T) {
			obj := Decode(t, `{
				"@context": "https://www.w3.org/ns/activitystreams",
				"type": "Document",
				"name": "4Q Sales Forecast",
				"url": {
					"type": "Link",
					"href": "http://example.org/4q-sales-forecast.pdf",
					"mediaType": "application/pdf"
				}
			}`)

			links := obj.URLs()
			require.Equal(t, 1, len(links))
			assertLink(t, links[0], pdfURL, pdfType)
			assert.Equal(t, pdfURL, obj.String("url"))
		})
	})
}

func Decode(t *testing.T, input string) *apubencoding.Object {
	dec := &apubencoding.Decoder{}
	obj, err := dec.Decode(strings.NewReader(input))
	require.Nil(t, err)
	return obj
}
