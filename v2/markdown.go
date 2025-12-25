package v2

import (
	"errors"
	"io"

	"github.com/charmbracelet/glamour"
)

// MarkdownRenderer is a pipeline wrapper around Option[*glamour.TermRenderer].
type MarkdownRenderer struct{ Option[*glamour.TermRenderer] }

// Markdown constructs a Glamour terminal renderer with reasonable defaults.
func Markdown(opts ...glamour.TermRendererOption) MarkdownRenderer {
	if len(opts) == 0 {
		tr, err := glamour.NewTermRenderer(
			glamour.WithEmoji(),
			glamour.WithStyles(glamour.DarkStyleConfig),
		)
		return MarkdownRenderer{Wrap(tr, err)}
	}
	tr, err := glamour.NewTermRenderer(opts...)
	return MarkdownRenderer{Wrap(tr, err)}
}

// Render renders markdown string.
func (m MarkdownRenderer) Render(markdown string) Str {
	if m.err != nil {
		return Str{Err[string](m.err)}
	}
	out, err := m.v.Render(markdown)
	if err != nil {
		_ = m.v.Close()
		return Str{Err[string](err)}
	}
	return Str{Ok(out)}
}

// RenderReader renders markdown from r.
func (m MarkdownRenderer) RenderReader(r io.Reader) Str {
	if m.err != nil {
		return Str{Err[string](m.err)}
	}
	if r == nil {
		return Str{Err[string](errors.New("lambda/v2: nil reader"))}
	}
	b := ReadAll(r)
	if b.err != nil {
		return Str{Err[string](b.err)}
	}
	return m.Render(b.String().Must())
}
