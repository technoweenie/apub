package apub_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/technoweenie/apub"
)

func TestRecipients(t *testing.T) {
	obj := Parse(t, `{
		"bto": "http://example.com/bto",
		"bcc": ["http://example.com/bcc", "http://example.com/cc"],
		"cc": "http://example.com/cc",
		"to": ["http://example.com/to", "http://example.com/public"],
		"actor": {
			"type": "Person",
			"id": "http://example.com/author"
		},
		"audience": ["http://example.com/public", "http://example.com/audience"],
		"target": {
			"type": "Note",
			"id": "http://example.com/target",
			"bto": "http://example.com/target/bto",
			"bcc": ["http://example.com/object/bcc", "http://example.com/object/cc"],
			"cc": "http://example.com/object/cc",
			"to": ["http://example.com/object/to", "http://example.com/public"]
		},
		"inReplyTo": {
			"type": "Note",
			"actor": "http://example.com/inReplyTo/author"
		},
		"object": {
			"attributedTo": "http://example.com/object/author",
			"bto": "http://example.com/object/bto",
			"bcc": ["http://example.com/object/bcc", "http://example.com/object/cc"],
			"cc": "http://example.com/object/cc",
			"to": ["http://example.com/object/to", "http://example.com/public"],
			"audience": ["http://example.com/public", "http://example.com/object/audience"],
			"tag": [
				{
					"type": "Person",
					"id": "http://example.com/tag/person"
				},
				{
					"type": "Mention",
					"id": "http://example.com/tag/mention",
					"href": "http://example.com/mention"
				},
				{
					"type": "Hashtag",
					"id": "http://example.com/hashtag/test",
					"name": "#test"
				}
			],
			"inReplyTo": {
				"type": "Note",
				"id": "http://example.com/object/inReplyTo",
				"bto": "http://example.com/object/inReplyTo/bto",
				"bcc": ["http://example.com/object/inReplyTo/bcc", "http://example.com/object/cc"],
				"cc": "http://example.com/object/cc",
				"to": ["http://example.com/object/to", "http://example.com/public"],
				"inReplyTo": {
					"type": "Note",
					"id": "http://example.com/inReplyTo^2",
					"bto": "http://example.com/inReplyTo^2/bto"
				}
			}
		}
	}`)

	recipients := apub.Recipients(obj)
	assert.Equal(t, 19, len(recipients), recipients)

	t.Run("activity", func(t *testing.T) {
		assert.Contains(t, recipients, "http://example.com/bto")
		assert.Contains(t, recipients, "http://example.com/bcc")
		assert.Contains(t, recipients, "http://example.com/cc")
		assert.Contains(t, recipients, "http://example.com/to")
		assert.Contains(t, recipients, "http://example.com/public")
		assert.Contains(t, recipients, "http://example.com/author")
		assert.Contains(t, recipients, "http://example.com/audience")
	})

	t.Run("target", func(t *testing.T) {
		assert.Contains(t, recipients, "http://example.com/target/bto")
		assert.Contains(t, recipients, "http://example.com/object/bcc")
		assert.Contains(t, recipients, "http://example.com/object/cc")
		assert.Contains(t, recipients, "http://example.com/object/to")
	})

	t.Run("inReplyTo", func(t *testing.T) {
		assert.Contains(t, recipients, "http://example.com/inReplyTo/author")
	})

	t.Run("object", func(t *testing.T) {
		assert.Contains(t, recipients, "http://example.com/object/author")
		assert.Contains(t, recipients, "http://example.com/object/bto")
		assert.Contains(t, recipients, "http://example.com/object/audience")
		assert.Contains(t, recipients, "http://example.com/tag/person")
		assert.Contains(t, recipients, "http://example.com/mention")

		t.Run("inReplyTo", func(t *testing.T) {
			assert.Contains(t, recipients, "http://example.com/object/inReplyTo/bto")
			assert.Contains(t, recipients, "http://example.com/object/inReplyTo/bcc")
		})
	})
}
