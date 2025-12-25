package v2

import "testing"

func TestMarkdown_Render(t *testing.T) {
	out := Markdown().Render("# Hello").Must()
	if out == "" {
		t.Fatalf("expected non-empty output")
	}
}


