package op

import (
	"bytes"
	"fmt"
	"io"

	. "github.com/emicklei/melrose"
)

type Serial struct {
	Target []Sequenceable
}

func (a Serial) S() Sequence {
	n := []Note{}
	for _, each := range a.Target {
		each.S().NotesDo(func(each Note) {
			n = append(n, each)
		})
	}
	return BuildSequence(n)
}

// Storex is part of Storable
func (a Serial) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "serial(")
	for i, each := range a.Target {
		s, ok := each.(Storable)
		if !ok {
			return ""
		}
		fmt.Fprintf(&b, "%s", s.Storex())
		if i < len(a.Target)-1 {
			io.WriteString(&b, ",")
		}
	}
	fmt.Fprintf(&b, ")")
	return b.String()
}

// Replaced is part of Replaceable
func (a Serial) Replaced(from, to Sequenceable) Sequenceable {
	if IsIdenticalTo(a, from) {
		return to
	}
	return Serial{Target: replacedAll(a.Target, from, to)}
}
