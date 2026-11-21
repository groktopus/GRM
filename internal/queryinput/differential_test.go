package queryinput

import (
	"encoding/json"
	"os"
	"testing"
)

func TestSchemaDifferentialCases(t *testing.T) {
	b, err := os.ReadFile("testdata/schema-cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Input string
		Valid bool
	}
	if err = json.Unmarshal(b, &cases); err != nil {
		t.Fatal(err)
	}
	for _, c := range cases {
		check(t, c.Input, c.Valid)
	}
}
