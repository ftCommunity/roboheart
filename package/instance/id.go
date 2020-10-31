package instance

const NON_INSTANCE_NAME = ""

type ID struct {
	Name     string //service name
	Instance string
}

func (id ID) equals(oid ID) bool {
	return id.Name == oid.Name && id.Instance == oid.Instance
}
