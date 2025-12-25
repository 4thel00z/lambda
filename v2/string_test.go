package v2

import "testing"

func TestLinesAndForEachLine(t *testing.T) {
	lines := StrOf("a\nb\nc\n").Lines().Must()
	if len(lines) != 3 || lines[0] != "a" || lines[2] != "c" {
		t.Fatalf("unexpected lines: %#v", lines)
	}

	upper := LinesOf(lines).ForEachLine(func(s string) string { return s + "!" }).Must()
	if len(upper) != 3 || upper[0] != "a!" || upper[2] != "c!" {
		t.Fatalf("unexpected mapped lines: %#v", upper)
	}
}
