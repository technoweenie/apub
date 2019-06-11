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
	s, _ := o.Fetch(key)
	return s
}

func (o *Object) Fetch(key string) (string, error) {
	ival, ok := o.data[key]
	if !ok {
		return "", nil
	}

	switch val := ival.(type) {
	case string:
		return val, nil
	case int:
		return strconv.Itoa(val), nil
	case int64:
		return strconv.FormatInt(val, 10), nil
	case map[string]interface{}:
		o2, err := o.valueAsObject(key, ival)
		if err != nil {
			return "", err
		}
		return o2.DefaultValue(), nil
	case []interface{}:
		objs, err := o.valueAsList(key, val)
		if err != nil {
			return "", err
		}
		if len(objs) == 0 {
			return "", nil
		}
		return objs[0].DefaultValue(), nil
	default:
		return fmt.Sprintf("%v", ival), nil
	}
}

func (o *Object) Object(key string) *Object {
	obj, _ := o.FetchObject(key)
	if obj == nil {
		return &Object{lang: o.lang, data: make(map[string]interface{})}
	}
	return obj
}

func (o *Object) FetchObject(key string) (*Object, error) {
	ival, ok := o.data[key]
	if !ok {
		return nil, nil
	}

	if list, ok := ival.([]interface{}); ok {
		objs, err := o.valueAsList(key, list)
		if err != nil {
			return nil, err
		}
		if len(objs) == 0 {
			return nil, nil
		}
		return objs[0], nil
	}

	return o.valueAsObject(key, ival)
}

func (o *Object) List(key string) []*Object {
	list, _ := o.FetchList(key)
	return list
}

func (o *Object) FetchList(key string) ([]*Object, error) {
	ival, ok := o.data[key]
	if !ok {
		return nil, nil
	}

	if list, ok := ival.([]interface{}); ok {
		return o.valueAsList(key, list)
	}

	obj, err := o.valueAsObject(key, ival)
	if obj == nil {
		return nil, err
	}
	return []*Object{obj}, err
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
