package op

import (
	"bytes"
	"fmt"

	"github.com/emicklei/melrose"
)

type Octave struct {
	Target []melrose.Sequenceable
	Offset melrose.Valueable
}

func (o Octave) S() melrose.Sequence {
	return Join{Target: o.Target}.S().Octaved(melrose.Int(o.Offset))
}

func (o Octave) Storex() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "octave(%v", o.Offset)
	appendStorexList(&b, false, o.Target)
	fmt.Fprintf(&b, ")")
	return b.String()
}

// Replaced is part of Replaceable
func (o Octave) Replaced(from, to melrose.Sequenceable) melrose.Sequenceable {
	if melrose.IsIdenticalTo(o, from) {
		return to
	}
	return Octave{Target: replacedAll(o.Target, from, to), Offset: o.Offset}
}
