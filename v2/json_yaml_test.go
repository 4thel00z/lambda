package v2

import "testing"

type testSpell struct {
	Name  string `json:"name" yaml:"name"`
	Power int    `json:"power" yaml:"power"`
}

func TestJSONRoundTrip(t *testing.T) {
	in := testSpell{Name: "bolt", Power: 7}
	b := Ok(in).ToJSON()
	if b.IsErr() {
		t.Fatalf("marshal err: %v", b.Err())
	}
	out := FromJSON[testSpell](b).Must()
	if out != in {
		t.Fatalf("mismatch: %#v vs %#v", out, in)
	}
}

func TestYAMLRoundTrip(t *testing.T) {
	in := testSpell{Name: "bolt", Power: 7}
	b := Ok(in).ToYAML()
	if b.IsErr() {
		t.Fatalf("marshal err: %v", b.Err())
	}
	out := FromYAML[testSpell](b).Must()
	if out != in {
		t.Fatalf("mismatch: %#v vs %#v", out, in)
	}
}
