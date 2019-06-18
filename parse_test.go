package apub_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/technoweenie/apub"
)

func TestParseObject(t *testing.T) {
	t.Run("basics", func(t *testing.T) {
		obj := Parse(t, `{
			"@context": "https://www.w3.org/ns/activitystreams",
			"type": "Object",
			"id": "http://www.test.example/object/1",
			"name": "A Simple, non-specific object",
			"num": 101,
			"roundup": 101.5,
			"rounddown": 103.45,
			"strnum": "104"
		}`)

		assert.Equal(t, "https://www.w3.org/ns/activitystreams", obj.Str("@context"))

		assert.Equal(t, "Object", obj.Str("type"))
		assert.Equal(t, obj.Str("type"), obj.Type())

		assert.Equal(t, "http://www.test.example/object/1", obj.Str("id"))
		assert.Equal(t, obj.Str("id"), obj.ID())

		assert.Equal(t, "A Simple, non-specific object", obj.Str("name"))
		assert.Equal(t, obj.Str("name"), obj.Name(""))

		assert.Equal(t, 101, obj.Int("num"))
		assert.Equal(t, 102, obj.Int("roundup"))
		assert.Equal(t, 103, obj.Int("rounddown"))
		assert.Equal(t, 104, obj.Int("strnum"))

		assert.Equal(t, "", obj.Str("not-a-property"))
		notObj := obj.Object("not-a-property")
		require.NotNil(t, notObj)
		assert.Equal(t, "", notObj.Type())
		notList := obj.List("not-a-list")
		assert.Equal(t, 0, len(notList))

		assert.Nil(t, obj.Errors())
		assert.NotNil(t, obj.NonFatalErrors())
	})

	t.Run("object url", func(t *testing.T) {
		assertLink := func(t *testing.T, link *apub.Object, expHref, expMediaType string) {
			require.NotNil(t, link)
			assert.Equal(t, "Link", link.Str("type"))
			assert.Equal(t, "Link", link.Type())
			assert.Equal(t, expHref, link.Str("href"))
			assert.Equal(t, expMediaType, link.Str("mediaType"))
		}

		pdfURL := "http://example.org/4q-sales-forecast.pdf"
		pdfType := "application/pdf"

		t.Run("string", func(t *testing.T) {
			obj := Parse(t, `{
				"@context": "https://www.w3.org/ns/activitystreams",
				"type": "Object",
				"name": "4Q Sales Forecast",
				"url": "http://example.org/4q-sales-forecast.pdf"
			}`)

			links := obj.URLs()
			require.Equal(t, 1, len(links))
			assertLink(t, links[0], pdfURL, "")
			assert.Equal(t, pdfURL, obj.Str("url"))

			assert.Nil(t, obj.Errors())
			assert.Nil(t, obj.NonFatalErrors())
		})

		t.Run("subclass string", func(t *testing.T) {
			obj := Parse(t, `{
				"@context": "https://www.w3.org/ns/activitystreams",
				"type": "Document",
				"name": "4Q Sales Forecast",
				"url": "http://example.org/4q-sales-forecast.pdf"
			}`)

			links := obj.URLs()
			require.Equal(t, 1, len(links))
			assertLink(t, links[0], pdfURL, "")
			assert.Equal(t, pdfURL, obj.Str("url"))

			assert.Nil(t, obj.Errors())
			assert.Nil(t, obj.NonFatalErrors())
		})

		t.Run("subclass object", func(t *testing.T) {
			obj := Parse(t, `{
				"@context": "https://www.w3.org/ns/activitystreams",
				"type": "Document",
				"name": "4Q Sales Forecast",
				"url": {
					"type": "Link",
					"href": "http://example.org/4q-sales-forecast.pdf",
					"mediaType": "application/pdf"
				},
				"icon": {
					"type": "Image",
					"url": [{
						"type": "Link",
						"href": "http://example.com/icon.jpg"
					}]
				},
				"image": [{
					"type": "Image",
					"url": "http://example.com/image.jpg"
				}]
			}`)

			t.Run("url", func(t *testing.T) {
				links := obj.URLs()
				require.Equal(t, 1, len(links))
				assertLink(t, links[0], pdfURL, pdfType)
				assert.Equal(t, pdfURL, obj.Str("url"))
			})

			t.Run("icon", func(t *testing.T) {
				assert.Equal(t, "http://example.com/icon.jpg", obj.Str("icon"))
				icons := obj.Icons()
				if assert.Equal(t, 1, len(icons)) {
					assert.Equal(t, "Image", icons[0].Type())
					assert.Equal(t, "http://example.com/icon.jpg", icons[0].Str("url"))
					urls := icons[0].URLs()
					if assert.Equal(t, 1, len(urls)) {
						assert.Equal(t, "Link", urls[0].Type())
						assert.Equal(t, "http://example.com/icon.jpg", urls[0].Str("href"))
					}
				}
			})

			t.Run("image", func(t *testing.T) {
				assert.Equal(t, "http://example.com/image.jpg", obj.Str("image"))
				images := obj.Images()
				if assert.Equal(t, 1, len(images)) {
					assert.Equal(t, "Image", images[0].Type())
					assert.Equal(t, "http://example.com/image.jpg", images[0].Str("url"))
					urls := images[0].URLs()
					if assert.Equal(t, 1, len(urls)) {
						assert.Equal(t, "Link", urls[0].Type())
						assert.Equal(t, "http://example.com/image.jpg", urls[0].Str("href"))
					}
				}
			})

			assert.Nil(t, obj.Errors())
			assert.Nil(t, obj.NonFatalErrors())
		})

		t.Run("subclass object", func(t *testing.T) {
			obj := Parse(t, `{
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

			assert.Nil(t, obj.Errors())
			assert.Nil(t, obj.NonFatalErrors())
		})
	})
}

func Parse(t *testing.T, input string) *apub.Object {
	dec := &apub.Parser{}
	obj, err := dec.Parse(strings.NewReader(input))
	require.Nil(t, err)
	return obj
}
