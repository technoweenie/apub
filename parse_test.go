package apub_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/technoweenie/apub"
	"golang.org/x/xerrors"
)

func TestParseObject(t *testing.T) {
	t.Run("basics", func(t *testing.T) {
		obj := Parse(t, `{
			"@context": "https://www.w3.org/ns/activitystreams",
			"type": "Object",
			"id": "http://www.test.example/object/1",
			"name": "A Simple, non-specific object",
			"true": true,
			"false": false,
			"t": "t",
			"f": "f",
			"t1": 1,
			"t1s": "1",
			"f0": 0,
			"f0s": "0",
			"invalidbool": "invalid",
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

		assert.Equal(t, "1", obj.Str("t1"))
		assert.Equal(t, "0", obj.Str("f0"))
		assert.Equal(t, "103.45", obj.Str("rounddown"))
		assert.Equal(t, "true", obj.Str("true"))
		assert.Equal(t, "false", obj.Str("false"))

		assert.Equal(t, 101, obj.Int("num"))
		assert.Equal(t, 102, obj.Int("roundup"))
		assert.Equal(t, 103, obj.Int("rounddown"))
		assert.Equal(t, 104, obj.Int("strnum"))

		assert.True(t, obj.Bool("true"))
		assert.True(t, obj.Bool("t"))
		assert.True(t, obj.Bool("t1"))
		assert.True(t, obj.Bool("t1s"))
		assert.False(t, obj.Bool("false"))
		assert.False(t, obj.Bool("f"))
		assert.False(t, obj.Bool("f0"))
		assert.False(t, obj.Bool("f0s"))

		assert.Equal(t, "", obj.Str("not-a-property"))
		notObj := obj.Object("not-a-property")
		require.NotNil(t, notObj)
		assert.Equal(t, "", notObj.Type())
		notList := obj.List("not-a-list")
		assert.Equal(t, 0, len(notList))

		assert.Nil(t, obj.Errors())
		nfErrs := obj.NonFatalErrors()
		if assert.Equal(t, 1, len(nfErrs), nfErrs) {
			assert.True(t, xerrors.Is(nfErrs[0], apub.ErrLangMapNotFound), nfErrs[0])
		}
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
