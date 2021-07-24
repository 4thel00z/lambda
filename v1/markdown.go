package v1

import (
	"errors"
	"github.com/charmbracelet/glamour"
	"io"
	"strings"
)

func Markdown(ops ...glamour.TermRendererOption) (o Option) {
	if len(ops) == 0 {
		return Wrap(glamour.NewTermRenderer(
			glamour.WithEmoji(),
			glamour.WithStyles(glamour.DarkStyleConfig),
		))
	}
	return Wrap(glamour.NewTermRenderer(ops...))
}

func MarkdownBlack(ops ...glamour.TermRendererOption) (o Option) {
	ops = append(ops, glamour.WithEmoji())
	ops = append(ops, glamour.WithStyles(glamour.DarkStyleConfig))
	return Wrap(glamour.NewTermRenderer(
		ops...,
	))
}

func (o Option) RenderFromReader(i io.ReadCloser) Option {

	renderer, ok := o.value.(*glamour.TermRenderer)
	if !ok {
		return Wrap(o.value, errors.New("o.value is not of type *glamour.TermRenderer"))
	}

	md, err := renderer.Render(Slurp(i).UnwrapString())
	if err != nil {
		return Wrap(md, renderer.Close())
	}

	return Wrap(md, err)
}

func (o Option) Render(i string) Option {
	return o.RenderFromReader(io.NopCloser(strings.NewReader(i)))
}
