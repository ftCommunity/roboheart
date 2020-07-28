package mapcompare

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/thoas/go-funk"
)

func TestCompareStringStringMap(t *testing.T) {
	for i := 0; i < 100; i++ {
		m0, m1, o0t, o1t, bst, bdt, e := genTestData(100)
		if e != nil {
			t.Error(e)
		}
		o0, o1, bs, bd := CompareStringStringMap(m0, m1)
		if !(checkStringSlicesEqual(o0t, o0) && checkStringSlicesEqual(o1t, o1) && checkStringSlicesEqual(bst, bs) && checkStringSlicesEqual(bdt, bd)) {
			t.Error(errors.New("Mismatch"))
		}
	}
}

func genTestData(n int) (m0, m1 map[string]string, o0, o1, bs, bd []string, e error) {
	m0 = map[string]string{}
	m1 = map[string]string{}
	o0 = []string{}
	o1 = []string{}
	bs = []string{}
	bd = []string{}
	var k, v, v0, v1 string
	for i := 0; i < n; i++ {
		k, v, e = genPair()
		if e != nil {
			return
		}
		m0[k] = v
		m1[k] = v
		bs = append(bs, k)
	}
	for i := 0; i < n; i++ {
		k, v0, v1, e = genTriple()
		if e != nil {
			return
		}
		m0[k] = v0
		m1[k] = v1
		bd = append(bd, k)
	}
	for i := 0; i < n; i++ {
		k, v, e = genPair()
		if e != nil {
			return
		}
		m0[k] = v
		o0 = append(o0, k)
	}
	for i := 0; i < n; i++ {
		k, v, e = genPair()
		if e != nil {
			return
		}
		m1[k] = v
		o1 = append(o1, k)
	}
	return
}

func genRandom() (s string, e error) {
	var u uuid.UUID
	u, e = uuid.NewRandom()
	if e != nil {
		return
	}
	s = u.String()
	return
}

func genPair() (s0, s1 string, e error) {
	s0, e = genRandom()
	if e != nil {
		return
	}
	s1, e = genRandom()
	return
}

func genTriple() (s0, s1, s2 string, e error) {
	s0, s1, e = genPair()
	if e != nil {
		return
	}
	s2, e = genRandom()
	return
}

func checkStringSlicesEqual(s0, s1 []string) bool {
	d0, d1 := funk.DifferenceString(s0, s1)
	return len(d0) == 0 && len(d1) == 0
}
