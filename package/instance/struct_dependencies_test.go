package instance

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDependencies_Compare(t *testing.T) {
	comp := func(d1, d2 Dependencies) bool {
		if len(d1) != len(d2) {
			return false
		}
	loop1:
		for _, id := range d1 {
			for _, oid := range d2 {
				if id.equals(oid) {
					continue loop1
				}
			}
			return false
		}
	loop2:
		for _, id := range d2 {
			for _, oid := range d1 {
				if id.equals(oid) {
					continue loop2
				}
			}
			return false
		}
		return true
	}
	test := func(t *testing.T) {
		td := make(Dependencies, 0)
		od := make(Dependencies, 0)
		ot := make(Dependencies, 0)
		oo := make(Dependencies, 0)
		b := make(Dependencies, 0)
		for i := 0; i < 10; i++ {
			tID := randID()
			oID := randID()
			bID := randID()

			td = append(td, tID)
			ot = append(ot, tID)

			od = append(od, oID)
			oo = append(oo, oID)

			td = append(td, bID)
			od = append(od, bID)
			b = append(b, bID)
		}
		ot1, oo1, b1 := td.Compare(od)
		oo2, ot2, b2 := od.Compare(td)
		assert.True(t, comp(ot1, ot2))
		assert.True(t, comp(oo1, oo2))
		assert.True(t, comp(b1, b2))
		assert.True(t, comp(ot1, ot))
		assert.True(t, comp(oo1, oo))
		assert.True(t, comp(b1, b))
	}
	for i := 0; i < 10; i++ {
		test(t)
	}
}
