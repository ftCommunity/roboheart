package relver

type ReleaseVersion interface {
	Update(token string) (error, bool)
	GetRelease() (Release, error)
	GetPreRelease() (Release, error)
	GetReleases() []Release
}
