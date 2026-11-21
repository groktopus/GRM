// Package queryinput validates the closed GRM QueryInput boundary.
package queryinput

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var ErrInput = errors.New("invalid input")

// Query is a validated QueryInput. Cursor and Question distinguish its closed branches.
type Query struct {
	Question         string
	Cursor           string
	Bindings         []Binding
	ValidAt, KnownAt *Instant
	PageSize         *int
}
type Binding struct{ Mention, Entity string }

func invalid() error { return ErrInput }

// scan rejects duplicate decoded keys at every object level before map decoding.
func scan(d *json.Decoder, depth int) error {
	if depth > 32 {
		return invalid()
	}
	t, err := d.Token()
	if err != nil {
		return err
	}
	delim, ok := t.(json.Delim)
	if !ok {
		return nil
	}
	switch delim {
	case '{':
		seen := map[string]bool{}
		for d.More() {
			key, err := d.Token()
			if err != nil {
				return err
			}
			k, ok := key.(string)
			if !ok || seen[k] {
				return invalid()
			}
			seen[k] = true
			if err := scan(d, depth+1); err != nil {
				return err
			}
		}
	case '[':
		for d.More() {
			if err := scan(d, depth+1); err != nil {
				return err
			}
		}
	default:
		return invalid()
	}
	_, err = d.Token()
	return err
}
func object(b []byte, allowed ...string) (map[string]json.RawMessage, error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(b, &m); err != nil || m == nil {
		return nil, invalid()
	}
	for k, v := range m {
		found := false
		for _, a := range allowed {
			if a == k {
				found = true
			}
		}
		if !found || bytes.Equal(bytes.TrimSpace(v), []byte("null")) {
			return nil, invalid()
		}
	}
	return m, nil
}
func str(b json.RawMessage, min, max int, pattern string) (string, error) {
	var s string
	if b == nil || json.Unmarshal(b, &s) != nil {
		return "", invalid()
	}
	n := utf8.RuneCountInString(s)
	if n < min || n > max || (pattern != "" && !regexp.MustCompile(pattern).MatchString(s)) {
		return "", invalid()
	}
	return s, nil
}

// Instant preserves the exact UTC seconds representation, including leap seconds.
// Values produced by ParseQuery can be ordered lexically and JSON encoded without
// normalizing an inserted second into the following midnight.
type Instant string

// Announced positive UTC leap boundaries through IERS Bulletin C72 (July 2026), Leap_Second.dat,
// https://hpiers.obspm.fr/iers/bul/bulc/Leap_Second.dat. Unknown
// boundaries are rejected; additions require reviewed evidence and regression cases.
const leapDates = "|1972-06-30|1972-12-31|1973-12-31|1974-12-31|1975-12-31|1976-12-31|1977-12-31|1978-12-31|1979-12-31|1981-06-30|1982-06-30|1983-06-30|1985-06-30|1987-12-31|1989-12-31|1990-12-31|1992-06-30|1993-06-30|1994-06-30|1995-12-31|1997-06-30|1998-12-31|2005-12-31|2008-12-31|2012-06-30|2015-06-30|2016-12-31|"

func instant(b json.RawMessage) (*Instant, error) {
	s, err := str(b, 20, 20, `^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z$`)
	if err != nil {
		return nil, err
	}
	calendar := s
	if s[17:19] == "60" {
		if s[11:] != "23:59:60Z" || !strings.Contains(leapDates, "|"+s[:10]+"|") {
			return nil, invalid()
		}
		calendar = s[:17] + "59Z"
	}
	if _, err := time.Parse("2006-01-02T15:04:05Z", calendar); err != nil {
		return nil, invalid()
	}
	value := Instant(s)
	return &value, nil
}

// ParseQuery validates the merged QueryInput shape, not authentication or retrieval semantics.
// The 65536-byte limit is the merged request bound; no source text is logged.
func ParseQuery(b []byte) (Query, error) {
	q, err := parseQuery(b)
	if err != nil {
		return Query{}, fmt.Errorf("query boundary: %w", ErrInput)
	}
	return q, nil
}
func parseQuery(b []byte) (Query, error) {
	var q Query
	if len(b) > 65536 || !utf8.Valid(b) {
		return q, invalid()
	}
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	if err := scan(d, 0); err != nil {
		return q, err
	}
	if _, err := d.Token(); err != io.EOF {
		return q, invalid()
	}
	m, err := object(b, "contract", "question", "bindings", "valid_at", "known_at", "page_size", "cursor")
	if err != nil {
		return q, err
	}
	c, err := str(m["contract"], 1, 1, "")
	if err != nil || c != "1" {
		return q, invalid()
	}
	if v, ok := m["cursor"]; ok {
		if len(m) != 2 {
			return q, invalid()
		}
		q.Cursor, err = str(v, 1, 512, `^[A-Za-z0-9_-]+$`)
		return q, err
	}
	q.Question, err = str(m["question"], 1, 16000, "")
	if err != nil {
		return q, err
	}
	if v, ok := m["page_size"]; ok {
		var number json.Number
		if json.Unmarshal(v, &number) != nil || len(v) == 0 || v[0] == '"' {
			return q, invalid()
		}
		// Bound exponent work before big.Rat expansion. With a nonzero value in
		// [1,20] and at most 65536 input bytes, larger exponents cannot qualify.
		if i := strings.IndexAny(string(number), "eE"); i >= 0 {
			exponent, e := strconv.ParseInt(string(number)[i+1:], 10, 32)
			if e != nil || exponent < -65536 || exponent > 65536 {
				return q, invalid()
			}
		}
		r, ok := new(big.Rat).SetString(string(number))
		if !ok || !r.IsInt() || r.Cmp(big.NewRat(1, 1)) < 0 || r.Cmp(big.NewRat(20, 1)) > 0 {
			return q, invalid()
		}
		n := int(r.Num().Int64())
		q.PageSize = &n
	}
	for k, dst := range map[string]**Instant{"valid_at": &q.ValidAt, "known_at": &q.KnownAt} {
		if v, ok := m[k]; ok {
			*dst, err = instant(v)
			if err != nil {
				return q, err
			}
		}
	}
	if v, ok := m["bindings"]; ok {
		var items []json.RawMessage
		if json.Unmarshal(v, &items) != nil || len(items) > 20 {
			return q, invalid()
		}
		q.Bindings = make([]Binding, 0, len(items))
		for _, item := range items {
			bm, err := object(item, "mention", "entity")
			if err != nil {
				return q, err
			}
			mention, err := str(bm["mention"], 1, 200, "")
			if err != nil {
				return q, err
			}
			entity, err := str(bm["entity"], 1, 70, `^ent_[A-Za-z0-9_-]{1,64}$`)
			if err != nil {
				return q, err
			}
			q.Bindings = append(q.Bindings, Binding{mention, entity})
		}
	}
	return q, nil
}
