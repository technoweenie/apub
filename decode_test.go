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

		assert.Equal(t, "https://www.w3.org/ns/activitystreams", obj.Str("@context"))
		assert.Equal(t, obj.Str("@context"), obj.Context())

		assert.Equal(t, "Object", obj.Str("type"))
		assert.Equal(t, obj.Str("type"), obj.Type())

		assert.Equal(t, "http://www.test.example/object/1", obj.Str("id"))
		assert.Equal(t, obj.Str("id"), obj.ID())

		assert.Equal(t, "A Simple, non-specific object", obj.Str("name"))
		assert.Equal(t, obj.Str("name"), obj.Name())

		assert.Equal(t, "", obj.Str("not-a-property"))
		notObj := obj.Object("not-a-property")
		require.NotNil(t, notObj)
		assert.Equal(t, "", notObj.Context())
		assert.Equal(t, "", notObj.Type())
		notList := obj.List("not-a-list")
		assert.Equal(t, 0, len(notList))
	})

	t.Run("object url", func(t *testing.T) {
		assertLink := func(t *testing.T, link *apubencoding.Object, expHref, expMediaType string) {
			require.NotNil(t, link)
			assert.Equal(t, "Link", link.Str("type"))
			assert.Equal(t, "Link", link.Type())
			assert.Equal(t, expHref, link.Str("href"))
			assert.Equal(t, expMediaType, link.Str("mediaType"))
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
			assert.Equal(t, pdfURL, obj.Str("url"))
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
			assert.Equal(t, pdfURL, obj.Str("url"))
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
			assert.Equal(t, pdfURL, obj.Str("url"))
		})

		t.Run("subclass object", func(t *testing.T) {
			obj := Decode(t, `{
				"@context": "https://www.w3.org/ns/activitystreams",
				"type": "Document",
				"name": "4Q Sales Forecast",
				"url": [
					{
						"type": "Link",
						"href": "http://example.org/4q-sales-forecast.pdf",
						"mediaType": "application/pdf"
					},
					{
						"type": "Link",
						"href": "http://example.org/4q-sales-forecast.html",
						"mediaType": "text/html"
					}
				]
			}`)

			links := obj.URLs()
			require.Equal(t, 2, len(links))
			assertLink(t, links[0], pdfURL, pdfType)
			assertLink(t, links[1], "http://example.org/4q-sales-forecast.html", "text/html")
			assert.Equal(t, pdfURL, obj.Str("url"))
		})
	})
}

func Decode(t *testing.T, input string) *apubencoding.Object {
	dec := &apubencoding.Decoder{}
	obj, err := dec.Decode(strings.NewReader(input))
	require.Nil(t, err)
	return obj
}
