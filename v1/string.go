package v1

import (
	"bufio"
	"log"
	"strings"
)

type StringOp func(string) string

func (o Option) ToStringLines() Option {
	if _, ok := o.value.([]string); ok {
		return o
	}
	f := o.UnwrapStringReader()
	lines := make([]string, 0)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	return Wrap(lines, sc.Err())
}
func (o Option) ForEachLine(fun StringOp) Option {
	if o.err != nil {
		log.Fatal(o.err)
	}

	if _, ok := o.value.([]string); !ok {
		log.Fatal("o.value is not of type []string")
	}

	lines := o.value.([]string)
	newLines := make([]string, len(lines))
	for _, line := range lines {
		newLines = append(newLines, fun(line))
	}
	return WrapValue(newLines)
}

func (o Option) ForEachLineReplace(m map[string]string) Option {
	return o.ForEachLine(ReplacerFromMap(m))
}

func (o Option) UnwrapStringLines() []string {
	res := o.ToStringLines()
	if res.err != nil {
		log.Fatal(res.err)
	}
	return res.Value().([]string)
}

func ReplacerFromMap(m map[string]string) StringOp {
	replacements := make([]string, 0, len(m)*2)
	for k, v := range m {
		replacements = append(replacements, k, v)
	}
	replacer := strings.NewReplacer(replacements...)
	return func(line string) string {
		return replacer.Replace(line)
	}
}
