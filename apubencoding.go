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
	return o.Str("@context")
}

func (o *Object) ID() string {
	return o.Str("id")
}

func (o *Object) Type() string {
	return o.Str("type")
}

func (o *Object) Name() string {
	return o.Str("name")
}

func (o *Object) URLs() []*Object {
	return o.List("url")
}

func (o *Object) Str(key string) string {
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
		o2, _ := o.valueAsObject(key, ival)
		return o2.DefaultValue()
	case []interface{}:
		objs, _ := o.valueAsList(key, val)
		if len(objs) == 0 {
			return ""
		}
		return objs[0].DefaultValue()
	default:
		return fmt.Sprintf("%v", ival)
	}
}

func (o *Object) Object(key string) *Object {
	ival, ok := o.data[key]
	if !ok {
		return nil
	}

	if list, ok := ival.([]interface{}); ok {
		objs, _ := o.valueAsList(key, list)
		if len(objs) == 0 {
			return nil
		}
		return objs[0]
	}

	obj, _ := o.valueAsObject(key, ival)
	return obj
}

func (o *Object) List(key string) []*Object {
	ival, ok := o.data[key]
	if !ok {
		return nil
	}

	if list, ok := ival.([]interface{}); ok {
		objs, _ := o.valueAsList(key, list)
		return objs
	}

	obj, _ := o.valueAsObject(key, ival)
	return []*Object{obj}
}

func (o *Object) valueAsObject(key string, ival interface{}) (*Object, error) {
	switch val := ival.(type) {
	case map[string]interface{}:
		return &Object{lang: o.lang, data: val}, nil
	case string:
		otype := o.Type()
		ptypes, ok := propertyTypes[otype]
		if !ok && otype == TypeObject {
			return nil, fmt.Errorf("unable to decode %s properties as objects", otype)
		}

		ptypes, ok = propertyTypes[TypeObject]
		if !ok {
			return nil, fmt.Errorf("unable to decode %s properties as objects", otype)
		}

		keyType, ok := ptypes[key]
		if !ok {
			return nil, fmt.Errorf("unable to decode %s %s property as object", otype, key)
		}

		return &Object{
			lang: o.lang,
			data: map[string]interface{}{
				"@context":        "https://www.w3.org/ns/activitystreams",
				"type":            keyType,
				defaults[keyType]: val,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unable to decode %T value as object: %+v", ival, ival)
	}
}

func (o *Object) DefaultValue() string {
	defkey, ok := defaults[o.Type()]
	if !ok {
		return o.Str("id")
	}
	return o.Str(defkey)
}

func (o *Object) valueAsList(key string, list []interface{}) ([]*Object, error) {
	objs := make([]*Object, 0, len(list))
	for _, iv := range list {
		o2, err := o.valueAsObject(key, iv)
		if err != nil {
			return objs, err
		}
		objs = append(objs, o2)
	}
	return objs, nil
}

const TypeObject = "Object"

var propertyTypes = map[string]map[string]string{
	TypeObject: map[string]string{
		"url": "Link",
	},
}

var defaults = map[string]string{
	"Link": "href",
}
