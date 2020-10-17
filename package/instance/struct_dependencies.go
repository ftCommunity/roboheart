package instance

type Dependencies []ID

func (d *Dependencies) Compare(od Dependencies) (ot, oo, b Dependencies) {
tloop:
	for _, id := range *d {
		for _, oid := range od {
			if id.equals(oid) {
				b = append(b, id)
				continue tloop
			}
		}
		ot = append(ot, id)
	}
oloop:
	for _, oid := range od {
		for _, id := range *d {
			if id.equals(oid) {
				continue oloop
			}
		}
		oo = append(oo, oid)
	}
	return
}

func (d *Dependencies) Add(id ID) {
	for _, oid := range *d {
		if id.equals(oid) {
			return
		}
	}
	*d = append(*d, id)
}

func (d *Dependencies) Delete(id ID) {
	n := Dependencies{}
	for _, oid := range *d {
		if !oid.equals(id) {
			n = append(n, oid)
		}
	}
	d = &n
}
