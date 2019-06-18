package apub_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/technoweenie/apub"
	"golang.org/x/xerrors"
)

func TestParseMastodon(t *testing.T) {
	t.Run("person", func(t *testing.T) {
		obj := Parse(t, `{
			"@context": [
				"https://www.w3.org/ns/activitystreams",
				"https://w3id.org/security/v1",
				{
					"manuallyApprovesFollowers": "as:manuallyApprovesFollowers",
					"toot": "http://joinmastodon.org/ns#",
					"featured": {
						"@id": "toot:featured",
						"@type": "@id"
					},
					"alsoKnownAs": {
						"@id": "as:alsoKnownAs",
						"@type": "@id"
					},
					"movedTo": {
						"@id": "as:movedTo",
						"@type": "@id"
					},
					"schema": "http://schema.org#",
					"PropertyValue": "schema:PropertyValue",
					"value": "schema:value",
					"Hashtag": "as:Hashtag",
					"Emoji": "toot:Emoji",
					"IdentityProof": "toot:IdentityProof",
					"focalPoint": {
						"@container": "@list",
						"@id": "toot:focalPoint"
					}
				}
			],
			"id": "https://mastodon.gamedev.place/users/bob",
			"type": "Person",
			"following": "https://mastodon.gamedev.place/users/bob/following",
			"followers": "https://mastodon.gamedev.place/users/bob/followers",
			"inbox": "https://mastodon.gamedev.place/users/bob/inbox",
			"outbox": "https://mastodon.gamedev.place/users/bob/outbox",
			"featured": "https://mastodon.gamedev.place/users/bob/collections/featured",
			"preferredUsername": "bob",
			"name": "Robert Tables",
			"summary": "<p>Bob</p>",
			"url": "https://mastodon.gamedev.place/@bob",
			"manuallyApprovesFollowers": false,
			"publicKey": {
				"id": "https://mastodon.gamedev.place/users/bob#main-key",
				"owner": "https://mastodon.gamedev.place/users/bob",
				"publicKeyPem": "-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAs+LxLJfjz+6Yf+1nh8rp\\na/ugMbp1geZFm2AsGZyIyB7CP/wRuzO9WmGkRQNpJmgEaYsiPN0l0ZcwoUtkXp41\\nZdUIOjuftLdNZAAaYFXzMEfmN3yE9LG5zOT9B3RSH/93psujPt0xUcurpN4L/III\\nwo9HawigZXPSY5J79Y4kDUOIpdw0o/36h0cZwAhrG+VHfAaHI5hShNW+6VzpWujP\\nzFI7eTtgJYLwE0PyJqLDqInbFINf4JaJqtvk7dLYeCQhPV8FrZMlsmrMVOY4TdAI\\nNKPu6QujEQjvguJy60//XYkH8stu5nlXKUR6GuY4s/Mo1mAb/bXo6lwMITAWPCD0\\noQIDAQAB\\n-----END PUBLIC KEY-----\n"
			},
			"tag": [],
			"attachment": [],
			"endpoints": {
				"sharedInbox": "https://mastodon.gamedev.place/inbox"
			},
			"icon": {
				"type": "Image",
				"mediaType": "image/jpeg",
				"url": "https://example.com/icon.jpg"
			},
			"image": {
				"type": "Image",
				"mediaType": "image/jpeg",
				"url": "https://example.com/image.jpg"
			}
		}`)

		assert.Equal(t, "https://mastodon.gamedev.place/users/bob", obj.ID())
		assert.Equal(t, "Person", obj.Type())
		assert.Equal(t, "https://mastodon.gamedev.place/users/bob/following", obj.Str("following"))
		assert.Equal(t, "<p>Bob</p>", obj.Str("summary"))

		pubKey := obj.Object("publicKey")
		if assert.NotNil(t, pubKey) {
			assert.Equal(t, "https://mastodon.gamedev.place/users/bob", pubKey.Str("owner"))
			pem := pubKey.Str("publicKeyPem")
			assert.True(t, strings.HasPrefix(pem, "-----BEGIN PUBLIC KEY-----"), pem)
		}

		// /cc subclass object tests
		assert.Equal(t, "https://example.com/icon.jpg", obj.Str("icon"))
		assert.Equal(t, "https://example.com/image.jpg", obj.Str("image"))

		assert.Nil(t, obj.Errors())
		assert.Nil(t, obj.NonFatalErrors())
	})

	t.Run("note collection", func(t *testing.T) {
		obj := Parse(t, `{
			"@context": "https://www.w3.org/ns/activitystreams",
			"id": "https://mastodon.gamedev.place/users/bob/collections/featured",
			"type": "OrderedCollection",
			"totalItems": 1,
			"orderedItems": [
				{
					"id": "https://mastodon.gamedev.place/users/bob/statuses/4815162342",
					"type": "Note",
					"summary": null,
					"inReplyTo": null,
					"published": "2019-04-14T17:19:09Z",
					"url": "https://mastodon.gamedev.place/@bob/4815162342",
					"attributedTo": "https://mastodon.gamedev.place/users/bob",
					"to": [
						"https://www.w3.org/ns/activitystreams#Public"
					],
					"cc": [
						"https://mastodon.gamedev.place/users/bob/followers"
					],
					"sensitive": false,
					"atomUri": "https://mastodon.gamedev.place/users/bob/statuses/4815162342",
					"inReplyToAtomUri": null,
					"conversation": "tag:mastodon.gamedev.place,2019-04-14:objectId=4815162342:objectType=Conversation",
					"content": "<p>Content</p>",
					"contentMap": {
						"en": "<p>EN Content</p>"
					},
					"attachment": [],
					"tag": [
						{
							"type": "Hashtag",
							"href": "https://mastodon.gamedev.place/tags/activitypub",
							"name": "#activitypub"
						}
					],
					"replies": {
						"id": "https://mastodon.gamedev.place/users/bob/statuses/4815162342/replies",
						"type": "Collection",
						"first": {
							"type": "CollectionPage",
							"partOf": "https://mastodon.gamedev.place/users/bob/statuses/4815162342/replies",
							"items": []
						}
					}
				}
			]
		}`)

		assert.Equal(t, "OrderedCollection", obj.Type())
		assert.Equal(t, "https://mastodon.gamedev.place/users/bob/collections/featured", obj.ID())
		assert.Equal(t, 1, obj.Int("totalItems"))

		items := obj.List("orderedItems")
		require.Equal(t, 1, len(items))
		item := items[0]
		assert.Equal(t, "https://mastodon.gamedev.place/users/bob/statuses/4815162342", item.ID())
		assert.Equal(t, "Note", item.Type())
		assert.Equal(t, "<p>EN Content</p>", item.Content(""))
		assert.Equal(t, "<p>EN Content</p>", item.Content("es"))
		assert.Equal(t, time.Date(2019, 4, 14, 17, 19, 9, 0, time.UTC), item.Time("published"))
		assert.False(t, item.Bool("sensitive"))

		urls := item.URLs()
		if assert.Equal(t, 1, len(urls), urls) {
			assert.Equal(t, "https://mastodon.gamedev.place/@bob/4815162342", urls[0].Str("href"))
		}

		replies := item.List("replies")
		if assert.Equal(t, 1, len(replies), replies) {
			assert.Equal(t, "https://mastodon.gamedev.place/users/bob/statuses/4815162342/replies", replies[0].ID())
			assert.Equal(t, "Collection", replies[0].Type())
		}

		assert.Nil(t, obj.Errors())

		nfErrs := obj.NonFatalErrors()
		if assert.Equal(t, 1, len(nfErrs), nfErrs) {
			assert.True(t, xerrors.Is(nfErrs[0], apub.ErrLangNotFound), nfErrs[0])
		}
	})

	t.Run("note", func(t *testing.T) {
		obj := Parse(t, `{
			"@context": [
				"https://www.w3.org/ns/activitystreams",
				{
					"ostatus": "http://ostatus.org#",
					"atomUri": "ostatus:atomUri",
					"inReplyToAtomUri": "ostatus:inReplyToAtomUri",
					"conversation": "ostatus:conversation",
					"sensitive": "as:sensitive",
					"Hashtag": "as:Hashtag",
					"toot": "http://joinmastodon.org/ns#",
					"Emoji": "toot:Emoji",
					"focalPoint": {
						"@container": "@list",
						"@id": "toot:focalPoint"
					},
					"blurhash": "toot:blurhash"
				}
			],
			"id": "https://mastodon.gamedev.place/users/bob/statuses/4815162342",
			"type": "Note",
			"summary": null,
			"inReplyTo": null,
			"published": "2019-06-13T04:46:37Z",
			"url": "https://mastodon.gamedev.place/@bob/4815162342",
			"attributedTo": "https://mastodon.gamedev.place/users/bob",
			"to": [
				"https://www.w3.org/ns/activitystreams#Public"
			],
			"cc": [
				"https://mastodon.gamedev.place/users/bob/followers"
			],
			"sensitive": false,
			"atomUri": "https://mastodon.gamedev.place/users/bob/statuses/4815162342",
			"inReplyToAtomUri": null,
			"conversation": "tag:mastodon.gamedev.place,2019-06-13:objectId=4815162342:objectType=Conversation",
			"content": "<p>Content</p>",
			"contentMap": {
				"en": "<p>Content EN</p>"
			},
			"attachment": [],
			"tag": [],
			"replies": {
				"id": "https://mastodon.gamedev.place/users/bob/statuses/4815162342/replies",
				"type": "Collection",
				"first": {
					"type": "CollectionPage",
					"partOf": "https://mastodon.gamedev.place/users/bob/statuses/4815162342/replies",
					"items": []
				}
			}
		}`)

		assert.Equal(t, "https://mastodon.gamedev.place/users/bob/statuses/4815162342", obj.ID())
		assert.Equal(t, "Note", obj.Type())
		assert.Equal(t, "<p>Content</p>", obj.Str("content"))
		assert.Equal(t, "<p>Content EN</p>", obj.Content(""))
		assert.Equal(t, time.Date(2019, 6, 13, 4, 46, 37, 0, time.UTC), obj.Time("published"))

		assert.Equal(t, "https://mastodon.gamedev.place/@bob/4815162342", obj.Str("url"))
		urls := obj.URLs()
		if assert.Equal(t, 1, len(urls)) {
			assert.Equal(t, "https://mastodon.gamedev.place/@bob/4815162342", urls[0].Str("href"))
		}

		assert.Nil(t, obj.Errors())
		assert.Nil(t, obj.NonFatalErrors())
	})
}
