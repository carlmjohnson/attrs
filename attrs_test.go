package attrs_test

import (
	"testing"

	"github.com/carlmjohnson/attrs"
)

func equalMap(t *testing.T, want, got map[string]interface{}) {
	t.Helper()
	if len(want) != len(got) {
		t.Fatalf("want: %v; got: %v", want, got)
	}
	for k, v := range want {
		if got[k] != v {
			t.Fatalf("want: %v; got: %v", want, got)
		}
	}
	for k, v := range got {
		if want[k] != v {
			t.Fatalf("want: %v; got: %v", want, got)
		}
	}
}

func TestToMap(t *testing.T) {
	tcases := map[string]struct {
		in   interface{}
		want map[string]interface{}
	}{
		"empty": {
			struct{}{},
			map[string]interface{}{},
		},
		"string": {
			struct {
				Hello string `tag:"hello"`
			}{"1"},
			map[string]interface{}{
				"hello": "1",
			},
		},
		"int": {
			struct {
				Hello int `tag:"hello"`
			}{1},
			map[string]interface{}{
				"hello": 1,
			},
		},
		"string+int": {
			struct {
				S string
				I int
			}{"1", 1},
			map[string]interface{}{
				"S": "1",
				"I": 1,
			},
		},
		"string+int+cache": {
			struct {
				S string
				I int
			}{"2", 2},
			map[string]interface{}{
				"S": "2",
				"I": 2,
			},
		},
	}
	for name, tc := range tcases {
		t.Run(name, func(t *testing.T) {
			m := attrs.ToMap(tc.in, "tag")
			equalMap(t, tc.want, m)
		})
	}
}

var m map[string]interface{}

var ss = []interface{}{
	&struct {
		Hello string `map:"whatever"`
	}{"world"},
	&struct {
		Hello     string `map:"value"`
		Something int
	}{"world", 1},
	&struct {
		Float float64
	}{1.5},
}

func BenchmarkToMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m = attrs.ToMap(&ss[i%len(ss)], "map")
	}
}

var errsink error

func BenchmarkFromMap(b *testing.B) {
	m := map[string]interface{}{
		"whatever": "one",
		"value":    "two",
		"Float":    3,
	}
	for i := 0; i < b.N; i++ {
		errsink = attrs.FromMap(m, &ss[i%len(ss)], "map")
	}
}

func TestFromMap(t *testing.T) {
	m := map[string]interface{}{
		"hello": "0xFF",
	}
	s := struct {
		Hello int `tag:"hello"`
	}{}
	err := attrs.FromMap(m, &s, "tag")
	if err != nil {
		t.Fatal(err)
	}
	if s.Hello != 255 {
		t.Fatalf("%v", s)
	}
	m = map[string]interface{}{
		"whatever": "one",
	}
	err = attrs.FromMap(m, &ss[0], "map")
	if err != nil {
		t.Fatal(err)
	}
	if (ss[0]).(*struct {
		Hello string `map:"whatever"`
	}).Hello != "one" {
		t.Fatalf("%v", ss[0])
	}
}
