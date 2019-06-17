package apubencoding_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeMastodon(t *testing.T) {
	t.Run("actor", func(t *testing.T) {
		obj := Decode(t, `{
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

		assert.Equal(t, "https://example.com/icon.jpg", obj.Str("icon"))
		icons := obj.Icons()
		if assert.Equal(t, 1, len(icons)) {
			assert.Equal(t, "Image", icons[0].Type())
			assert.Equal(t, "image/jpeg", icons[0].Str("mediaType"))
			assert.Equal(t, "https://example.com/icon.jpg", icons[0].Str("url"))
			urls := icons[0].URLs()
			if assert.Equal(t, 1, len(urls)) {
				assert.Equal(t, "https://example.com/icon.jpg", urls[0].Str("href"))
			}
		}

		assert.Equal(t, "https://example.com/image.jpg", obj.Str("image"))
		images := obj.Images()
		if assert.Equal(t, 1, len(images)) {
			assert.Equal(t, "Image", images[0].Type())
			assert.Equal(t, "image/jpeg", images[0].Str("mediaType"))
			assert.Equal(t, "https://example.com/image.jpg", images[0].Str("url"))
			urls := images[0].URLs()
			if assert.Equal(t, 1, len(urls)) {
				assert.Equal(t, "https://example.com/image.jpg", urls[0].Str("href"))
			}
		}

		assert.Nil(t, obj.Errors())
		assert.Nil(t, obj.NonFatalErrors())
	})

	t.Run("note", func(t *testing.T) {
		obj := Decode(t, `{
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

		assert.Equal(t, "https://mastodon.gamedev.place/@bob/4815162342", obj.Str("url"))
		urls := obj.URLs()
		if assert.Equal(t, 1, len(urls)) {
			assert.Equal(t, "https://mastodon.gamedev.place/@bob/4815162342", urls[0].Str("href"))
		}

		assert.Nil(t, obj.Errors())
		assert.Nil(t, obj.NonFatalErrors())
	})
}