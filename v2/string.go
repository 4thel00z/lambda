package v2

import (
	"bufio"
	"strings"
)

// StringOp transforms a single line.
type StringOp func(string) string

// Lines splits a string into lines (scanner-based).
func (s Str) Lines() Lines {
	if s.err != nil {
		return Lines{Err[[]string](s.err)}
	}

	sc := bufio.NewScanner(strings.NewReader(s.v))
	// Allow long lines (1MiB).
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	lines := make([]string, 0)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return Lines{Wrap(lines, sc.Err())}
}

// Lines splits bytes into lines by converting to string first.
func (b Bytes) Lines() Lines { return b.String().Lines() }

// ForEachLine applies fun to each line.
func (l Lines) ForEachLine(fun StringOp) Lines {
	if l.err != nil {
		return Lines{Err[[]string](l.err)}
	}
	if fun == nil {
		return Lines{Err[[]string](ErrNilFunc("ForEachLine"))}
	}
	lines := l.v
	newLines := make([]string, len(lines))
	for i, line := range lines {
		newLines[i] = fun(line)
	}
	return Lines{Ok(newLines)}
}

// ForEachLineReplace applies a replacer generated from m to each line.
func (l Lines) ForEachLineReplace(m map[string]string) Lines {
	return l.ForEachLine(ReplacerFromMap(m))
}

// ReplacerFromMap returns a StringOp that applies all key/value replacements.
func ReplacerFromMap(m map[string]string) StringOp {
	replacements := make([]string, 0, len(m)*2)
	for k, v := range m {
		replacements = append(replacements, k, v)
	}
	replacer := strings.NewReplacer(replacements...)
	return func(line string) string { return replacer.Replace(line) }
}
