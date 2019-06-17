package apubencoding_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/technoweenie/apubencoding"
	"golang.org/x/xerrors"
)

func TestDecode(t *testing.T) {
	t.Run("valueAsObject errors", func(t *testing.T) {
		obj := Decode(t, `{
			"type": "Object",
		  "object1": {
				"type": "TestObject",
				"icon": 123,
				"test": "test"
			}
		}`)

		obj1, err := obj.FetchObject("object1")
		assert.Nil(t, err)

		missing, err := obj1.FetchObject("missing")
		assert.Nil(t, missing)
		assert.Nil(t, err)

		_, err = obj1.FetchObject("test")
		assert.True(t, xerrors.Is(err, apubencoding.ErrKeyNotObject), err)

		_, err = obj1.FetchObject("icon")
		assert.True(t, xerrors.Is(err, apubencoding.ErrKeyTypeNotObject), err)
	})

	t.Run("lang errors", func(t *testing.T) {
		obj := Decode(t, `{
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
		assert.True(t, xerrors.Is(err, apubencoding.ErrLangMapNotFound))

		img := obj.Object("image")
		require.NotNil(t, img)
		imgName, err := img.FetchLang("name", "en")
		assert.Equal(t, "image", imgName)
		assert.Nil(t, err)

		esName, err := img.FetchLang("name", "es")
		assert.Equal(t, "image", esName)
		assert.True(t, xerrors.Is(err, apubencoding.ErrLangNotFound))
	})

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
		assert.Equal(t, obj.Str("name"), obj.Name(""))

		assert.Equal(t, "", obj.Str("not-a-property"))
		notObj := obj.Object("not-a-property")
		require.NotNil(t, notObj)
		assert.Equal(t, "", notObj.Context())
		assert.Equal(t, "", notObj.Type())
		notList := obj.List("not-a-list")
		assert.Equal(t, 0, len(notList))

		assert.Nil(t, obj.Errors())
		assert.NotNil(t, obj.NonFatalErrors())
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

			assert.Nil(t, obj.Errors())
			assert.Nil(t, obj.NonFatalErrors())
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

			assert.Nil(t, obj.Errors())
			assert.Nil(t, obj.NonFatalErrors())
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

			assert.Nil(t, obj.Errors())
			assert.Nil(t, obj.NonFatalErrors())
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

			assert.Nil(t, obj.Errors())
			assert.Nil(t, obj.NonFatalErrors())
		})
	})

	t.Run("content language map", func(t *testing.T) {
		obj := Decode(t, `{
			"@context": "https://www.w3.org/ns/activitystreams",
			"name": "Cane Sugar Processing",
			"type": "Note",
			"summaryMap": {
				"en": "A simple <em>note</em>",
				"es": "Una <em>nota</em> sencilla",
				"zh-Hans": "一段<em>简单的</em>笔记"
			},
			"contentMap": {
				"es": "Una <em>nota</em> sencilla"
			}
		}`)

		t.Run("summary", func(t *testing.T) {
			defaultLang, err := obj.FetchLang("summary", "")
			assert.Equal(t, "A simple <em>note</em>", defaultLang)
			assert.Equal(t, defaultLang, obj.Summary(""))
			assert.Nil(t, err)

			esLang, err := obj.FetchLang("summary", "es")
			assert.Equal(t, "Una <em>nota</em> sencilla", esLang)
			assert.Equal(t, esLang, obj.Summary("es"))
			assert.Nil(t, err)

			otherLang, err := obj.FetchLang("summary", "other")
			assert.Equal(t, "A simple <em>note</em>", otherLang)
			assert.Equal(t, otherLang, obj.Summary("other"))
			assert.NotNil(t, err)
			assert.False(t, apubencoding.FatalLangErr(err))
			assert.True(t, xerrors.Is(err, apubencoding.ErrLangNotFound), err)

			assert.Nil(t, obj.Errors())
			assert.NotNil(t, obj.NonFatalErrors())
		})

		t.Run("name", func(t *testing.T) {
			defaultLang, err := obj.FetchLang("name", "")
			assert.Equal(t, "Cane Sugar Processing", defaultLang)
			assert.Equal(t, defaultLang, obj.Name(""))
			assert.NotNil(t, err)
			assert.False(t, apubencoding.FatalLangErr(err))
			assert.True(t, xerrors.Is(err, apubencoding.ErrLangMapNotFound), err)

			esLang, err := obj.FetchLang("name", "es")
			assert.Equal(t, "Cane Sugar Processing", esLang)
			assert.Equal(t, esLang, obj.Name("es"))
			assert.NotNil(t, err)
			assert.False(t, apubencoding.FatalLangErr(err))
			assert.True(t, xerrors.Is(err, apubencoding.ErrLangMapNotFound), err)

			otherLang, err := obj.FetchLang("name", "other")
			assert.Equal(t, "Cane Sugar Processing", otherLang)
			assert.Equal(t, otherLang, obj.Name("other"))
			assert.NotNil(t, err)
			assert.False(t, apubencoding.FatalLangErr(err))
			assert.True(t, xerrors.Is(err, apubencoding.ErrLangMapNotFound), err)

			assert.Nil(t, obj.Errors())
			assert.NotNil(t, obj.NonFatalErrors())
		})

		t.Run("content", func(t *testing.T) {
			defaultLang, err := obj.FetchLang("content", "")
			assert.Equal(t, "", defaultLang)
			assert.Equal(t, defaultLang, obj.Content(""))
			assert.NotNil(t, err)
			assert.False(t, apubencoding.FatalLangErr(err))
			assert.True(t, xerrors.Is(err, apubencoding.ErrLangNotFound), err)

			esLang, err := obj.FetchLang("content", "es")
			assert.Equal(t, "Una <em>nota</em> sencilla", esLang)
			assert.Equal(t, esLang, obj.Content("es"))
			assert.Nil(t, err)

			otherLang, err := obj.FetchLang("content", "other")
			assert.Equal(t, "", otherLang)
			assert.Equal(t, otherLang, obj.Content("other"))
			assert.NotNil(t, err)
			assert.False(t, apubencoding.FatalLangErr(err))
			assert.True(t, xerrors.Is(err, apubencoding.ErrLangNotFound), err)

			assert.Nil(t, obj.Errors())
			assert.NotNil(t, obj.NonFatalErrors())
		})
	})
}

func Decode(t *testing.T, input string) *apubencoding.Object {
	dec := &apubencoding.Decoder{}
	obj, err := dec.Decode(strings.NewReader(input))
	require.Nil(t, err)
	return obj
}
