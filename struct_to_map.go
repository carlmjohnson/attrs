package attrs

import (
	"reflect"
	"sync"
)

var structmapCache sync.Map

type cacheKey struct {
	reflect.Type
	string
}

type fieldkey struct {
	pos int
	key string
}

func caches2mGet(s reflect.Value, tag string) []fieldkey {
	t := s.Type()
	v, ok := structmapCache.Load(cacheKey{t, tag})
	if ok {
		return v.([]fieldkey)
	}
	nf := s.NumField()
	fieldkeys := make([]fieldkey, 0, nf)
	for i := 0; i < nf; i++ {
		if f := s.Field(i); f.CanInterface() {
			key := t.Field(i).Name
			if tagkey := t.Field(i).Tag.Get(tag); tagkey != "" {
				key = tagkey
			}
			fieldkeys = append(fieldkeys, fieldkey{i, key})
		}
	}
	caches2mSet(t, tag, fieldkeys)
	return fieldkeys
}

func caches2mSet(t reflect.Type, tag string, fieldkeys []fieldkey) {
	structmapCache.Store(cacheKey{t, tag}, fieldkeys)
}

func ToMap(v interface{}, tag string) map[string]interface{} {
	s := reflect.ValueOf(v)
	for s.Kind() == reflect.Ptr || s.Kind() == reflect.Interface {
		s = s.Elem()
	}

	fieldkeys := caches2mGet(s, tag)

	r := make(map[string]interface{}, s.NumField())
	for _, k := range fieldkeys {
		r[k.key] = s.Field(k.pos).Interface()
	}
	return r
}
