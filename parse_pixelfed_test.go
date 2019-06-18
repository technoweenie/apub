package apub_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePixelfed(t *testing.T) {
	t.Run("person", func(t *testing.T) {
		obj := Parse(t, `{
			"@context": [
				"https://www.w3.org/ns/activitystreams",
				"https://w3id.org/security/v1",
				{
					"manuallyApprovesFollowers": "as:manuallyApprovesFollowers"
				}
			],
			"id": "https://fedi.pictures/users/Rob_T_Firefly",
			"type": "Person",
			"following": "https://fedi.pictures/users/bob/following",
			"followers": "https://fedi.pictures/users/bob/followers",
			"inbox": "https://fedi.pictures/users/bob/inbox",
			"outbox": "https://fedi.pictures/users/bob/outbox",
			"preferredUsername": "bob",
			"name": "bob",
			"summary": "what about bob",
			"url": "https://fedi.pictures/bob",
			"manuallyApprovesFollowers": false,
			"publicKey": {
				"id": "https://fedi.pictures/users/bob#main-key",
				"owner": "https://fedi.pictures/users/bob",
				"publicKeyPem": "-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAs+LxLJfjz+6Yf+1nh8rp\\na/ugMbp1geZFm2AsGZyIyB7CP/wRuzO9WmGkRQNpJmgEaYsiPN0l0ZcwoUtkXp41\\nZdUIOjuftLdNZAAaYFXzMEfmN3yE9LG5zOT9B3RSH/93psujPt0xUcurpN4L/III\\nwo9HawigZXPSY5J79Y4kDUOIpdw0o/36h0cZwAhrG+VHfAaHI5hShNW+6VzpWujP\\nzFI7eTtgJYLwE0PyJqLDqInbFINf4JaJqtvk7dLYeCQhPV8FrZMlsmrMVOY4TdAI\\nNKPu6QujEQjvguJy60//XYkH8stu5nlXKUR6GuY4s/Mo1mAb/bXo6lwMITAWPCD0\\noQIDAQAB\\n-----END PUBLIC KEY-----\n"
			},
			"icon": {
				"type": "Image",
				"mediaType": "image/jpeg",
				"url": "https://example.com/icon.jpg"
			}
		}`)

		assert.Equal(t, "https://fedi.pictures/users/Rob_T_Firefly", obj.ID())
		assert.Equal(t, "Person", obj.Type())
		assert.False(t, obj.Bool("manuallyApprovesFollowers"))
		assert.Equal(t, "https://example.com/icon.jpg", obj.Str("icon"))

		icons := obj.Icons()
		if assert.Equal(t, 1, len(icons), icons) {
			assert.Equal(t, "Image", icons[0].Type())
			assert.Equal(t, "image/jpeg", icons[0].Str("mediaType"))
			assert.Equal(t, "https://example.com/icon.jpg", icons[0].Str("url"))
		}

		assert.Nil(t, obj.Errors())
		assert.Nil(t, obj.NonFatalErrors())
	})

	t.Run("note", func(t *testing.T) {
		obj := Parse(t, `{
			"@context": [
				"https://www.w3.org/ns/activitystreams",
				"https://w3id.org/security/v1",
				{
					"sc": "http://schema.org#",
					"Hashtag": "as:Hashtag",
					"sensitive": "as:sensitive",
					"commentsEnabled": "sc:Boolean",
					"capabilities": {
						"announce": {
							"@type": "@id"
						},
						"like": {
							"@type": "@id"
						},
						"reply": {
							"@type": "@id"
						}
					}
				}
			],
			"id": "https://fedi.pictures/p/bob/4815162342",
			"type": "Note",
			"summary": null,
			"content": "content",
			"inReplyTo": null,
			"published": "2019-04-30T22:01:40+00:00",
			"url": "https://fedi.pictures/p/bob/4815162342",
			"attributedTo": "https://fedi.pictures/users/bob",
			"to": [
				"https://www.w3.org/ns/activitystreams#Public"
			],
			"cc": [
				"https://fedi.pictures/users/bob/followers"
			],
			"sensitive": false,
			"attachment": [
				{
					"type": "Image",
					"mediaType": "image/jpeg",
					"url": "https://example.com/image.jpg",
					"name": "some image"
				}
			],
			"tag": [],
			"commentsEnabled": true,
			"capabilities": {
				"announce": "https://www.w3.org/ns/activitystreams#Public",
				"like": "https://www.w3.org/ns/activitystreams#Public",
				"reply": "https://www.w3.org/ns/activitystreams#Public"
			}
		}`)

		assert.Equal(t, "https://fedi.pictures/p/bob/4815162342", obj.ID())
		assert.Equal(t, "Note", obj.Type())
		assert.Equal(t, "content", obj.Content(""))
		assert.False(t, obj.Bool("sensitive"))
		assert.True(t, obj.Bool("commentsEnabled"))
	})
}
