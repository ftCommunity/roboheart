package filesystem

import "github.com/spf13/afero"

type FileSystem interface {
	afero.Fs
}
