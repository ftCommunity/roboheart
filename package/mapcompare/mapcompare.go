package mapcompare

import (
	"github.com/thoas/go-funk"
)

func CompareStringStringMap(m0, m1 map[string]string) (o0, o1, bs, bd []string) {
	o0 = []string{}
	o1 = []string{}
	bs = []string{}
	bd = []string{}
	e := []string{}
	for k, v0 := range m0 {
		e = append(e, k)
		if v1, f := m1[k]; f {
			if v0 == v1 {
				bs = append(bs, k)
			} else {
				bd = append(bd, k)
			}
		} else {
			o0 = append(o0, k)
		}
	}
	for k, v1 := range m1 {
		if funk.ContainsString(e, k) {
			continue
		}
		if v0, f := m0[k]; f {
			if v0 == v1 {
				bs = append(bs, k)
			} else {
				bd = append(bd, k)
			}
		} else {
			o1 = append(o1, k)
		}
	}
	return
}
