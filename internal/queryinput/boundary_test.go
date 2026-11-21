package queryinput

import (
	"errors"
	"strings"
	"testing"
)

func TestQueryBoundary(t *testing.T) {
	cases := []struct {
		name, input string
		valid       bool
	}{
		{"fresh", `{"contract":"1","question":"Who?"}`, true},
		{"cursor", `{"contract":"1","cursor":"next_1"}`, true},
		{"full", `{"contract":"1","question":"誰？","bindings":[{"mention":"誰","entity":"ent_A"}],"valid_at":"2024-02-29T00:00:00Z","known_at":"2026-11-21T12:00:00Z","page_size":20}`, true},
		{"integer decimal", `{"contract":"1","question":"x","page_size":1.0}`, true},
		{"integer exponent", `{"contract":"1","question":"x","page_size":2e1}`, true},
		{"missing contract", `{"question":"x"}`, false},
		{"missing branch", `{"contract":"1"}`, false},
		{"null question", `{"contract":"1","question":null}`, false},
		{"null optional", `{"contract":"1","question":"x","bindings":null}`, false},
		{"both branches", `{"contract":"1","question":"x","cursor":"c"}`, false},
		{"cursor option", `{"contract":"1","cursor":"c","page_size":1}`, false},
		{"empty cursor", `{"contract":"1","cursor":""}`, false},
		{"bad token", `{"contract":"1","cursor":"a b"}`, false},
		{"unknown", `{"contract":"1","question":"x","extra":1}`, false},
		{"wrong case", `{"Contract":"1","question":"x"}`, false},
		{"duplicate", `{"contract":"1","question":"x","question":"y"}`, false},
		{"escaped duplicate", `{"contract":"1","question":"x","\u0071uestion":"y"}`, false},
		{"nested duplicate", `{"contract":"1","question":"x","bindings":[{"mention":"a","mention":"b","entity":"ent_A"}]}`, false},
		{"nested unknown", `{"contract":"1","question":"x","bindings":[{"mention":"a","entity":"ent_A","x":0}]}`, false},
		{"nested required", `{"contract":"1","question":"x","bindings":[{"mention":"a"}]}`, false},
		{"nested null", `{"contract":"1","question":"x","bindings":[null]}`, false},
		{"trailing value", `{"contract":"1","question":"x"} {}`, false},
		{"trailing junk", `{"contract":"1","question":"x"} x`, false},
		{"bad calendar", `{"contract":"1","question":"x","valid_at":"2025-02-29T00:00:00Z"}`, false},
		{"non UTC", `{"contract":"1","question":"x","known_at":"2026-01-01T00:00:00+00:00"}`, false},
		{"fractional seconds", `{"contract":"1","question":"x","valid_at":"2026-01-01T00:00:00.1Z"}`, false},
		{"fractional page", `{"contract":"1","question":"x","page_size":1.1}`, false},
		{"zero page", `{"contract":"1","question":"x","page_size":0}`, false},
		{"large page", `{"contract":"1","question":"x","page_size":21}`, false},
		{"string page", `{"contract":"1","question":"x","page_size":"2"}`, false},
		{"unicode limit", `{"contract":"1","question":"` + strings.Repeat("界", 16000) + `"}`, true},
		{"unicode excess", `{"contract":"1","question":"` + strings.Repeat("界", 16001) + `"}`, false},
		{"cursor excess", `{"contract":"1","cursor":"` + strings.Repeat("a", 513) + `"}`, false},
		{"request bound", strings.Repeat(" ", 65537), false},
		{"invalid UTF8", string([]byte{0xff}), false},
		{"huge exponent", `{"contract":"1","question":"x","page_size":1e999999999}`, false},
		{"binding limit", `{"contract":"1","question":"x","bindings":[` + strings.Repeat(`{"mention":"m","entity":"ent_A"},`, 20) + `{"mention":"m","entity":"ent_A"}]}`, false},
		{"array", `[]`, false},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseQuery([]byte(tt.input))
			if (err == nil) != tt.valid {
				t.Fatalf("valid=%v error=%v", tt.valid, err)
			}
			if err != nil && !errors.Is(err, ErrInput) {
				t.Fatal("input identity lost")
			}
		})
	}
}

func FuzzParseQuery(f *testing.F) {
	for _, s := range []string{`{"contract":"1","question":"x"}`, `{"contract":"1","cursor":"n"}`, `null`, `{"a":{"a":1,"a":2}}`} {
		f.Add([]byte(s))
	}
	f.Fuzz(func(t *testing.T, b []byte) {
		q, err := ParseQuery(b)
		if err == nil && ((q.Question == "") == (q.Cursor == "")) {
			t.Fatal("invalid union accepted")
		}
	})
}
