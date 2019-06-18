package apub_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/technoweenie/apub"
	"golang.org/x/xerrors"
)

func TestParseContent(t *testing.T) {
	obj := Parse(t, `{
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
		assert.False(t, apub.FatalLangErr(err))
		assert.True(t, xerrors.Is(err, apub.ErrLangNotFound), err)

		assert.Nil(t, obj.Errors())
		assert.NotNil(t, obj.NonFatalErrors())
	})

	t.Run("name", func(t *testing.T) {
		defaultLang, err := obj.FetchLang("name", "")
		assert.Equal(t, "Cane Sugar Processing", defaultLang)
		assert.Equal(t, defaultLang, obj.Name(""))
		assert.NotNil(t, err)
		assert.False(t, apub.FatalLangErr(err))
		assert.True(t, xerrors.Is(err, apub.ErrLangMapNotFound), err)

		esLang, err := obj.FetchLang("name", "es")
		assert.Equal(t, "Cane Sugar Processing", esLang)
		assert.Equal(t, esLang, obj.Name("es"))
		assert.NotNil(t, err)
		assert.False(t, apub.FatalLangErr(err))
		assert.True(t, xerrors.Is(err, apub.ErrLangMapNotFound), err)

		otherLang, err := obj.FetchLang("name", "other")
		assert.Equal(t, "Cane Sugar Processing", otherLang)
		assert.Equal(t, otherLang, obj.Name("other"))
		assert.NotNil(t, err)
		assert.False(t, apub.FatalLangErr(err))
		assert.True(t, xerrors.Is(err, apub.ErrLangMapNotFound), err)

		assert.Nil(t, obj.Errors())
		assert.NotNil(t, obj.NonFatalErrors())
	})

	t.Run("content", func(t *testing.T) {
		defaultLang, err := obj.FetchLang("content", "")
		assert.Equal(t, "", defaultLang)
		assert.Equal(t, defaultLang, obj.Content(""))
		assert.NotNil(t, err)
		assert.False(t, apub.FatalLangErr(err))
		assert.True(t, xerrors.Is(err, apub.ErrLangNotFound), err)

		esLang, err := obj.FetchLang("content", "es")
		assert.Equal(t, "Una <em>nota</em> sencilla", esLang)
		assert.Equal(t, esLang, obj.Content("es"))
		assert.Nil(t, err)

		otherLang, err := obj.FetchLang("content", "other")
		assert.Equal(t, "", otherLang)
		assert.Equal(t, otherLang, obj.Content("other"))
		assert.NotNil(t, err)
		assert.False(t, apub.FatalLangErr(err))
		assert.True(t, xerrors.Is(err, apub.ErrLangNotFound), err)

		assert.Nil(t, obj.Errors())
		assert.NotNil(t, obj.NonFatalErrors())
	})
}
