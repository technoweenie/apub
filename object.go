package apub

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

var DefaultLang = "en"

type Object struct {
	path         []string
	lang         string
	data         map[string]interface{}
	errors       []error
	nonFatal     []error
	addError     func(error)
	addLangError func(error)
}

func New(data map[string]interface{}) *Object {
	obj := &Object{lang: DefaultLang, data: data}
	if ty := obj.Type(); len(ty) > 0 {
		obj.path = []string{ty}
	} else {
		obj.path = []string{"UnknownType"}
	}

	obj.addError = func(err error) {
		obj.errors = append(obj.errors, err)
	}
	obj.addLangError = func(err error) {
		if FatalLangErr(err) {
			obj.errors = append(obj.errors, err)
			return
		}
		obj.nonFatal = append(obj.nonFatal, err)
	}
	return obj
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
		o.addError(err)
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
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
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

func (o *Object) Int(key string) int {
	i, err := o.FetchInt(key)
	if err != nil {
		o.addError(err)
	}
	return i
}

func (o *Object) FetchInt(key string) (int, error) {
	ival, ok := o.data[key]
	if !ok {
		return 0, nil
	}

	switch val := ival.(type) {
	case float64:
		return int(math.Round(val)), nil
	case string:
		i, err := strconv.Atoi(val)
		if err != nil {
			return i, xerrors.Errorf("FetchInt: %q: %w", val, ErrInvalidInt)
		}
		return i, nil
	default:
		return 0, xerrors.Errorf("FetchInt: %T %+v: %w", ival, ival, ErrInvalidInt)
	}
}

func (o *Object) Bool(key string) bool {
	b, err := o.FetchBool(key)
	if err != nil {
		o.addError(err)
	}
	return b
}

func (o *Object) FetchBool(key string) (bool, error) {
	ival, ok := o.data[key]
	if !ok {
		return false, nil
	}

	switch val := ival.(type) {
	case bool:
		return val, nil
	case string:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return b, xerrors.Errorf("FetchBool: %q: %w", val, err)
		}
		return b, nil
	case float64:
		if n := int64(math.Round(val)); n == 1 {
			return true, nil
		}
		return false, nil
	default:
		return false, xerrors.Errorf("FetchBool: %T %+v: %w", ival, ival, ErrInvalidBool)
	}
}

func (o *Object) Object(key string) *Object {
	obj, err := o.FetchObject(key)
	if err != nil {
		o.addError(err)
	}
	if obj == nil {
		return o.newObj(key, make(map[string]interface{}))
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
		o.addError(err)
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
		o.addLangError(err)
	}
	return s
}

func (o *Object) Name(lang string) string {
	s, err := o.FetchLang("name", lang)
	if err != nil {
		o.addLangError(err)
	}
	return s
}

func (o *Object) Summary(lang string) string {
	s, err := o.FetchLang("summary", lang)
	if err != nil {
		o.addLangError(err)
	}
	return s
}

func (o *Object) FetchLang(key, lang string) (string, error) {
	if len(lang) == 0 {
		lang = o.lang
	}
	if len(lang) == 0 {
		return o.Str(key), xerrors.Errorf("FetchLang: %s.%s in %q: %w",
			strings.Join(o.path, "."), key, lang, ErrLangNotFound)
	}

	cmap, err := o.FetchObject(key + "Map")
	if err != nil || cmap == nil {
		if err == nil {
			err = ErrLangMapNotFound
		}
		return o.Str(key), xerrors.Errorf("FetchLang: %s.%s: %w",
			strings.Join(o.path, "."), key, err)
	}

	val, err := cmap.Fetch(lang)
	if err != nil || len(val) == 0 {
		if err == nil {
			err = ErrLangNotFound
		}
		noLangErr := xerrors.Errorf("FetchLang: %s.%s in %q: %w",
			strings.Join(o.path, "."), key, lang, err)
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

func (o *Object) valueAsObject(key string, ival interface{}) (*Object, error) {
	switch val := ival.(type) {
	case map[string]interface{}:
		return o.newObj(key, val), nil
	case string:
		otype := o.Type()
		var ptypes map[string]string

		if pt, ok := propertyTypes[otype]; ok {
			ptypes = pt
		} else {
			ptypes = propertyTypes[TypeObject] // always exists
		}

		keyType, ok := ptypes[key]
		if !ok {
			return nil, xerrors.Errorf("valueAsObject: (%s) %s key %q: %w",
				otype, strings.Join(o.path, "."), key, ErrKeyNotObject)
		}

		return o.newObj(key, map[string]interface{}{
			"type":            keyType,
			defaults[keyType]: val,
		}), nil
	default:
		return nil, xerrors.Errorf("valueAsObject: (%s) %s key %q: (%T) %+v: %w",
			o.Type(), strings.Join(o.path, "."), key, ival, ival, ErrKeyTypeNotObject)
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

func (o *Object) newObj(key string, data map[string]interface{}) *Object {
	return &Object{
		path:         append(o.path, key),
		lang:         o.lang,
		data:         data,
		addError:     o.addError,
		addLangError: o.addLangError,
	}
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
