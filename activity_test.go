package apub_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/technoweenie/apub"
)

func TestActivityCreate(t *testing.T) {
	obj := Parse(t, `{
		"@context": "https://www.w3.org/ns/activitystreams",
		"type": "Note",
		"content": "This is a note",
		"published": "2015-02-10T15:04:55Z",
		"to": ["https://example.org/~john/"],
		"cc": ["https://example.com/~erik/followers",
					 "https://www.w3.org/ns/activitystreams#Public"]
	}`)

	assert.Equal(t, "https://www.w3.org/ns/activitystreams", obj.Str("@context"))
	assert.Equal(t, "Note", obj.Type())
	assert.Equal(t, "This is a note", obj.Content(""))
	assert.Equal(t, time.Date(2015, 2, 10, 15, 4, 55, 0, time.UTC), obj.Time("published"))
	assert.Equal(t, []string{"https://example.org/~john/"}, obj.To())
	assert.Equal(t, []string{"https://example.com/~erik/followers",
		"https://www.w3.org/ns/activitystreams#Public"}, obj.CC())

	act := apub.CreateActivity(obj)

	assert.Equal(t, "https://www.w3.org/ns/activitystreams", act.Str("@context"))
	assert.Equal(t, "Create", act.Type())
	assert.Equal(t, time.Date(2015, 2, 10, 15, 4, 55, 0, time.UTC), act.Time("published"))
	assert.Equal(t, []string{"https://example.org/~john/"}, act.To())
	assert.Equal(t, []string{"https://example.com/~erik/followers",
		"https://www.w3.org/ns/activitystreams#Public"}, act.CC())

	note := act.Object("object")
	assert.Equal(t, "", note.Str("@context"))
	assert.Equal(t, "Note", note.Type())
	assert.Equal(t, "This is a note", note.Content(""))
	assert.Equal(t, time.Date(2015, 2, 10, 15, 4, 55, 0, time.UTC), note.Time("published"))
	assert.Equal(t, []string{"https://example.org/~john/"}, note.To())
	assert.Equal(t, []string{"https://example.com/~erik/followers",
		"https://www.w3.org/ns/activitystreams#Public"}, note.CC())
}
