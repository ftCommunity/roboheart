package filesystem

import "github.com/spf13/afero"

func NewMock() FileSystem {
	return afero.NewMemMapFs()
}
