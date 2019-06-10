package apubencoding

import (
	"fmt"
	"strconv"
)

const DefaultLang = "en"

type Object struct {
	lang string
	data map[string]interface{}
}

func (o *Object) Context() string {
	return o.String("@context")
}

func (o *Object) ID() string {
	return o.String("id")
}

func (o *Object) Type() string {
	return o.String("type")
}

func (o *Object) Name() string {
	return o.String("name")
}

func (o *Object) URL() *Object {
	return o.Object("url")
}

func (o *Object) String(key string) string {
	ival, ok := o.data[key]
	if !ok {
		return ""
	}

	switch val := ival.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case map[string]interface{}:
		o2 := &Object{lang: o.lang, data: val}
		defkey, ok := defaults[o2.Type()]
		if !ok {
			return fmt.Sprintf("%v", val)
		}
		return o2.String(defkey)
	default:
		return fmt.Sprintf("%v", ival)
	}
}

func (o *Object) Object(key string) *Object {
	ival, ok := o.data[key]
	if !ok {
		return nil
	}

	switch val := ival.(type) {
	case map[string]interface{}:
		return &Object{lang: o.lang, data: val}
	case string:
		otype := o.Type()
		ptypes, ok := propertyTypes[otype]
		if !ok {
			return nil
		}

		keyType, ok := ptypes[key]
		if !ok {
			return nil
		}

		return &Object{
			lang: o.lang,
			data: map[string]interface{}{
				"@context":        "https://www.w3.org/ns/activitystreams",
				"type":            keyType,
				defaults[keyType]: val,
			},
		}
	default:
		fmt.Printf("welp! %T %+v\n", ival, ival)
		return nil
	}
}

var propertyTypes = map[string]map[string]string{
	"Object": map[string]string{
		"url": "Link",
	},
}

var defaults = map[string]string{
	"Link": "href",
}
