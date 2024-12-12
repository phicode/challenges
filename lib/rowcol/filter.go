package rowcol

import "github.com/phicode/challenges/lib/assert"

func MinPos(ps []Pos) Pos {
	l := len(ps)
	assert.True(l > 0)
	m := ps[0]
	for i := 1; i < l; i++ {
		p := ps[i]
		if p.Row < m.Row || (p.Row == m.Row && p.Col < m.Col) {
			m = p
		}
	}
	return m
}
