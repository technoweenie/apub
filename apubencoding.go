package apubencoding

import (
	"errors"
	"fmt"
	"strconv"

	"golang.org/x/xerrors"
)

const DefaultLang = "en"

type Object struct {
	lang     string
	data     map[string]interface{}
	errors   []error
	nonFatal []error
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

func (o *Object) Icons() []*Object {
	return o.List("icon")
}

func (o *Object) Images() []*Object {
	return o.List("image")
}

func (o *Object) URLs() []*Object {
	return o.List("url")
}

func (o *Object) Str(key string) string {
	s, err := o.Fetch(key)
	if err != nil {
		o.errors = append(o.errors, err)
	}
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
	obj, err := o.FetchObject(key)
	if err != nil {
		o.errors = append(o.errors, err)
	}
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
	list, err := o.FetchList(key)
	if err != nil {
		o.errors = append(o.errors, err)
	}
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

func (o *Object) Content(lang string) string {
	s, err := o.FetchLang("content", lang)
	if err != nil {
		if FatalLangErr(err) {
			o.errors = append(o.errors, err)
		} else {
			o.nonFatal = append(o.nonFatal, err)
		}
	}
	return s
}

func (o *Object) Name(lang string) string {
	s, err := o.FetchLang("name", lang)
	if err != nil {
		if FatalLangErr(err) {
			o.errors = append(o.errors, err)
		} else {
			o.nonFatal = append(o.nonFatal, err)
		}
	}
	return s
}

func (o *Object) Summary(lang string) string {
	s, err := o.FetchLang("summary", lang)
	if err != nil {
		if FatalLangErr(err) {
			o.errors = append(o.errors, err)
		} else {
			o.nonFatal = append(o.nonFatal, err)
		}
	}
	return s
}

func (o *Object) FetchLang(key, lang string) (string, error) {
	if len(lang) == 0 {
		lang = o.lang
	}
	if len(lang) == 0 {
		return o.Str(key), xerrors.Errorf("FetchLang: %q in %q: %w", key, lang, ErrLangNotFound)
	}

	cmap, err := o.FetchObject(key + "Map")
	if err != nil || cmap == nil {
		if err == nil {
			err = ErrLangMapNotFound
		}
		return o.Str(key), xerrors.Errorf("FetchLang: %q: %w", key, err)
	}

	val, err := cmap.Fetch(lang)
	if err != nil || len(val) == 0 {
		if err == nil {
			err = ErrLangNotFound
		}
		noLangErr := xerrors.Errorf("FetchLang: %q in %q: %w", key, lang, err)
		if lang != o.lang {
			fallback, _ := o.FetchLang(key, o.lang)
			return fallback, noLangErr
		}
		return o.Str(key), noLangErr
	}

	return val, nil
}

func (o *Object) Errors() []error {
	return o.errors
}

func (o *Object) NonFatalErrors() []error {
	return o.nonFatal
}

var ErrLangNotFound = errors.New("key not translated to given language")
var ErrLangMapNotFound = errors.New("key has no language map")

func FatalLangErr(err error) bool {
	if xerrors.Is(err, ErrLangNotFound) {
		return false
	}
	if xerrors.Is(err, ErrLangMapNotFound) {
		return false
	}
	return true
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
		"icon":  "Image",
		"image": "Image",
		"url":   "Link",
	},
}

var defaults = map[string]string{
	"Image": "url",
	"Link":  "href",
}
