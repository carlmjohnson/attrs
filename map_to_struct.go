package attrs

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
)

var mapstructcache sync.Map

type fieldkeytype struct {
	pos int
	key string
	reflect.Type
}

func cachem2sGet(s reflect.Value, tag string) []fieldkeytype {
	t := s.Type()
	v, ok := mapstructcache.Load(cacheKey{t, tag})
	if ok {
		return v.([]fieldkeytype)
	}
	nf := s.NumField()
	fieldkeytypes := make([]fieldkeytype, 0, nf)
	for i := 0; i < nf; i++ {
		if f := s.Field(i); f.CanSet() {
			tf := t.Field(i)
			key := tf.Name
			if tagkey := tf.Tag.Get(tag); tagkey != "" {
				key = tagkey
			}
			fieldkeytypes = append(fieldkeytypes, fieldkeytype{i, key, tf.Type})
		}
	}
	cachem2sSet(t, tag, fieldkeytypes)
	return fieldkeytypes
}

func cachem2sSet(t reflect.Type, tag string, fieldkeytypes []fieldkeytype) {
	mapstructcache.Store(cacheKey{t, tag}, fieldkeytypes)
}

func FromMap(m map[string]interface{}, v interface{}, tag string) error {
	s := reflect.ValueOf(v)
	for s.Kind() == reflect.Ptr || s.Kind() == reflect.Interface {
		s = s.Elem()
	}
	fieldkeytypes := cachem2sGet(s, tag)
	if len(fieldkeytypes) == 0 {
		panic("can't set any fields on v")
	}
	for _, k := range fieldkeytypes {
		key := k.key
		if m[key] == nil {
			continue
		}
		mv := reflect.ValueOf(m[key])
		f := s.Field(k.pos)
		switch {
		case mv.Type().AssignableTo(k.Type):
			f.Set(mv)
		case k.Kind() == reflect.String && mv.Kind() == reflect.Int:
			f.SetString(strconv.FormatInt(mv.Int(), 10))
		case k.Kind() == reflect.Int && mv.Kind() == reflect.String:
			n, err := strconv.ParseInt(mv.String(), 0, 0)
			if err != nil {
				return fmt.Errorf("field %s could not convert %T to %s: %v",
					key, m[key], k.Type.String(), err)
			}
			f.SetInt(n)
		case mv.Type().ConvertibleTo(k.Type):
			f.Set(mv.Convert(k.Type))
		default:
			return fmt.Errorf("field %s could not convert %T to %s",
				key, m[key], k.Type.String())
		}
	}

	return nil
}
