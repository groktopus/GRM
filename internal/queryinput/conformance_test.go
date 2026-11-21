package queryinput

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func check(t *testing.T, input string, want bool) Query {
	t.Helper()
	q, err := ParseQuery([]byte(input))
	if (err == nil) != want {
		t.Fatalf("accept=%v want=%v input=%s", err == nil, want, input)
	}
	if err != nil && (!errors.Is(err, ErrInput) || !reflect.DeepEqual(q, Query{})) {
		t.Fatal("failure exposed partial value or lost identity")
	}
	return q
}
func encoded(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}
func TestPreservation(t *testing.T) {
	q := check(t, `{"contract":"1","question":" e\u0301\n😀 ","bindings":[{"mention":" first\n","entity":"ent_A"},{"mention":"😀","entity":"ent_B"}],"page_size":20,"valid_at":"2016-12-31T23:59:60Z","known_at":"2017-01-01T00:00:00Z"}`, true)
	if q.Question != " e\u0301\n😀 " || !reflect.DeepEqual(q.Bindings, []Binding{{" first\n", "ent_A"}, {"😀", "ent_B"}}) || q.PageSize == nil || *q.PageSize != 20 || q.ValidAt == nil || *q.ValidAt != "2016-12-31T23:59:60Z" || q.KnownAt == nil || *q.KnownAt != "2017-01-01T00:00:00Z" {
		t.Fatalf("lost fields: %+v", q)
	}
	if encoded(q.ValidAt) != `"2016-12-31T23:59:60Z"` {
		t.Fatal("leap normalization")
	}
	minimal := check(t, `{"contract":"1","question":"q"}`, true)
	if minimal.Bindings != nil || minimal.PageSize != nil || minimal.ValidAt != nil || minimal.KnownAt != nil {
		t.Fatal("fabricated fields")
	}
	empty := check(t, `{"contract":"1","question":"q","bindings":[]}`, true)
	if empty.Bindings == nil {
		t.Fatal("lost explicit empty collection")
	}
	cursor := check(t, `{"contract":"1","cursor":"next_-1"}`, true)
	if cursor.Cursor != "next_-1" || cursor.Question != "" {
		t.Fatal("cursor lost")
	}
}
func TestUnicode(t *testing.T) {
	for _, tt := range []struct{ raw, want string }{
		{`"😀"`, "😀"}, {`"\ud83d\ude00"`, "😀"}, {`"�"`, "�"}, {`"\ufffd"`, "�"}, {`"é"`, "é"}, {`"\u00e9"`, "é"}, {`"e\u0301"`, "e\u0301"}, {`"\\ud800"`, `\ud800`},
		// Characterization of Go replacement, not a lossless GRM guarantee.
		{`"\ud800"`, "�"}, {`"\udc00"`, "�"}, {`"\ud800A"`, "�A"}, {`"\udc00\ud800"`, "��"},
	} {
		q := check(t, `{"contract":"1","question":`+tt.raw+`,"bindings":[{"mention":`+tt.raw+`,"entity":"ent_A"}]}`, true)
		if q.Question != tt.want || q.Bindings[0].Mention != tt.want {
			t.Fatalf("unicode %s => %+v", tt.raw, q)
		}
	}
}
func TestDateOracle(t *testing.T) {
	good := []string{"2016-12-31T23:59:59Z", "2016-12-31T23:59:60Z", "2017-01-01T00:00:00Z", "2000-02-29T00:00:00Z", "1972-06-30T23:59:60Z"}
	bad := []string{"2016-12-31T23:58:60Z", "2016-12-30T23:59:60Z", "2017-01-01T23:59:60Z", "2016-12-31T22:59:60Z", "2016-12-31T23:59:61Z", "2016-12-31T24:00:00Z", "1900-02-29T00:00:00Z", "2016-02-30T00:00:00Z", "2016-13-01T00:00:00Z", "2016-12-31T23:59:60.0Z", "2016-12-31T23:59:60+00:00", "2017-01-01T00:59:60+01:00", "2016-12-31T15:59:60-08:00", "2016-12-31t23:59:60z", "2016-12-31", "", "2030-06-30T23:59:60Z"}
	for _, field := range []string{"valid_at", "known_at"} {
		for _, s := range good {
			q := check(t, `{"contract":"1","question":"q","`+field+`":`+encoded(s)+`}`, true)
			v := q.ValidAt
			if field == "known_at" {
				v = q.KnownAt
			}
			if v == nil || string(*v) != s {
				t.Fatal("date lost")
			}
		}
		for _, s := range bad {
			check(t, `{"contract":"1","question":"q","`+field+`":`+encoded(s)+`}`, false)
		}
	}
	if !(Instant(good[0]) < Instant(good[1]) && Instant(good[1]) < Instant(good[2])) {
		t.Fatal("adjacent ordering")
	}
}
func TestRequiredAndTypes(t *testing.T) {
	wrong := []any{nil, false, 1, []any{}, map[string]any{}}
	for _, branch := range []string{"question", "cursor"} {
		for _, field := range []string{"contract", branch} {
			for _, v := range append(wrong, "") {
				m := map[string]any{"contract": "1", branch: "q"}
				m[field] = v
				check(t, encoded(m), false)
			}
			m := map[string]any{"contract": "1", branch: "q"}
			delete(m, field)
			check(t, encoded(m), false)
		}
	}
	for _, field := range []string{"mention", "entity"} {
		for _, v := range append(wrong, "") {
			b := map[string]any{"mention": "m", "entity": "ent_A"}
			b[field] = v
			check(t, encoded(map[string]any{"contract": "1", "question": "q", "bindings": []any{b}}), false)
		}
		b := map[string]any{"mention": "m", "entity": "ent_A"}
		delete(b, field)
		check(t, encoded(map[string]any{"contract": "1", "question": "q", "bindings": []any{b}}), false)
	}
	for _, field := range []string{"bindings", "valid_at", "known_at", "page_size"} {
		for _, v := range []any{nil, false, "x", map[string]any{}, []any{}} {
			if field == "bindings" && reflect.DeepEqual(v, []any{}) {
				continue
			}
			check(t, encoded(map[string]any{"contract": "1", "question": "q", field: v}), false)
		}
		check(t, encoded(map[string]any{"contract": "1", "cursor": "q", field: []any{}}), false)
	}
	check(t, `{"contract":"2","question":"q"}`, false)
	for _, field := range []string{"bindings", "valid_at", "known_at"} {
		check(t, encoded(map[string]any{"contract": "1", "question": "q", field: 1}), false)
	}
	for _, s := range []string{`{}`, `null`, `1`, `true`, `"x"`, `[]`, `{"contract":"1","question":"q","bindings":[{}]}`} {
		check(t, s, false)
	}
}
func TestEndpoints(t *testing.T) {
	for _, tt := range []struct {
		field        string
		min, max     int
		prefix, unit string
	}{
		{"question", 1, 16000, "", "😀"}, {"cursor", 1, 512, "", "a"}, {"mention", 1, 200, "", "😀"}, {"entity", 1, 64, "ent_", "a"},
	} {
		for _, n := range []int{0, tt.min, tt.max, tt.max + 1} {
			value := tt.prefix + strings.Repeat(tt.unit, n)
			m := map[string]any{"contract": "1", "question": "q"}
			if tt.field == "mention" || tt.field == "entity" {
				b := map[string]any{"mention": "m", "entity": "ent_A"}
				b[tt.field] = value
				m["bindings"] = []any{b}
			} else {
				m[tt.field] = value
				if tt.field == "cursor" {
					delete(m, "question")
				}
			}
			check(t, encoded(m), n >= tt.min && n <= tt.max)
		}
	}
	for _, n := range []int{0, 20, 21} {
		b := make([]any, n)
		for i := range b {
			b[i] = map[string]any{"mention": "m", "entity": "ent_A"}
		}
		check(t, encoded(map[string]any{"contract": "1", "question": "q", "bindings": b}), n <= 20)
	}
	for _, tt := range []struct {
		number string
		want   bool
	}{{"1", true}, {"20", true}, {"1.0", true}, {"1e0", true}, {"0", false}, {"21", false}, {"1.5", false}, {"true", false}, {`"1"`, false}, {"null", false}, {"[]", false}, {"{}", false}} {
		check(t, `{"contract":"1","question":"q","page_size":`+tt.number+`}`, tt.want)
	}
	for _, value := range []string{"ent_", "ent_a\n", "ent_é", "wrong_A"} {
		check(t, encoded(map[string]any{"contract": "1", "question": "q", "bindings": []any{map[string]any{"mention": "m", "entity": value}}}), false)
	}
}

// TestYearZeroCharacterization records the GRM/Go outcome separately from the
// pinned Python checker, whose date-time checker rejects year zero. That checker
// limitation must not become a universal RFC3339/schema rejection oracle.
func TestYearZeroCharacterization(t *testing.T) {
	const value = "0000-01-01T00:00:00Z"
	for _, field := range []string{"valid_at", "known_at"} {
		q := check(t, encoded(map[string]any{"contract": "1", "question": "q", field: value}), true)
		got := q.ValidAt
		if field == "known_at" {
			got = q.KnownAt
		}
		if got == nil || string(*got) != value || encoded(got) != encoded(value) {
			t.Fatalf("%s year-zero preservation: %+v", field, q)
		}
	}
}
