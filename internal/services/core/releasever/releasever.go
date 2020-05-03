package relver

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/blang/semver"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/acm"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
	"github.com/ftCommunity/roboheart/package/threadmanager"
	"github.com/google/go-github/v31/github"
)

const (
	PERMISSION_BASE   = "releasever"
	PERMISSION_UPDATE = PERMISSION_BASE + "." + "update"
)

type relver struct {
	release, prerelease *Release
	releases            []Release
	lock                sync.Mutex
	last                time.Time
	gh                  *github.Client
	logger              service.LoggerFunc
	error               service.ErrorFunc
	tm                  *threadmanager.ThreadManager
	acm                 acm.ACM
}

type ReleaseVersion interface {
	Update(token string) error
	GetRelease() (Release, error)
	GetPreRelease() (Release, error)
	GetReleases() []Release
}

type Release struct {
	Version  semver.Version
	Download string
}

func (r *relver) Init(services map[string]service.Service, logger service.LoggerFunc, e service.ErrorFunc) error {
	r.logger = logger
	r.error = e
	if err := servicehelpers.CheckMainDependencies(r, services); err != nil {
		return err
	}
	var ok bool
	r.acm, ok = services["acm"].(acm.ACM)
	if !ok {
		return errors.New("Type assertion error")
	}
	r.acm.RegisterPermission(PERMISSION_UPDATE, map[string]bool{"root": true, "user": true, "app": false})
	r.gh = github.NewClient(nil)
	r.tm = threadmanager.NewThreadManager(r.logger, r.error)
	r.tm.Load("update", r.updateThread)
	r.tm.Start("update")
	return nil
}

func (r *relver) Stop() error                                                { return nil }
func (r *relver) Name() string                                               { return "relver" }
func (r *relver) Dependencies() ([]string, []string)                         { return []string{"acm"}, []string{} }
func (r *relver) SetAdditionalDependencies(map[string]service.Service) error { return nil }
func (r *relver) UnsetAdditionalDependencies(s chan interface{})             { s <- struct{}{} }

func (r *relver) updateThread(logger service.LoggerFunc, e service.ErrorFunc, stop, stopped chan interface{}) {
	if err := r.getReleaseData(); err != nil {
		logger(err)
	}
	for {
		select {
		case <-stop:
			{
				stopped <- struct{}{}
				return
			}
		case <-time.After(15 * time.Minute):
			{
				if err := r.getReleaseData(); err != nil {
					logger(err)
				}
			}
		}
	}
}

func (r *relver) getReleaseData() error {
	releases, _, err := r.gh.Repositories.ListReleases(context.Background(), "ftCommunity", "ftcommunity-TXT", nil)
	if err != nil {
		return err
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	r.release, r.prerelease = nil, nil
	r.releases = make([]Release, 0)
	for _, rrel := range releases {
		if rrel.TagName != nil {
			rel, err := newRelease(rrel)
			if err != nil {
				return err
			}
			r.releases = append(r.releases, *rel)
			pr := false
			if rrel.Prerelease != nil {
				pr = *rrel.Prerelease
			}
			if r.release == nil && !pr {
				r.release = rel
			}
			if r.prerelease == nil && pr {
				r.prerelease = rel
			}
		}
	}
	return nil
}

func (r *relver) Update(token string) error {
	if err := acm.CheckTokenPermission(r.acm, token, PERMISSION_UPDATE); err != nil {
		return err
	}
	return r.getReleaseData()
}

func (r *relver) GetRelease() (Release, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.release == nil {
		return Release{}, errors.New("Prerelease not set")
	}
	return *r.release, nil
}

func (r *relver) GetPreRelease() (Release, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.prerelease == nil {
		return Release{}, errors.New("Prerelease not set")
	}
	return *r.prerelease, nil
}

func (r *relver) GetReleases() []Release {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.releases
}

func newRelease(rel *github.RepositoryRelease) (*Release, error) {
	r := new(Release)
	ver, err := semver.Make(strings.Replace(*rel.TagName, "v", "", -1))
	if err != nil {
		return nil, err
	}
	r.Version = ver
	if r.Download, err = getAssetURL(rel); err != nil {
		return nil, err
	}
	return r, nil
}

func getAssetURL(rel *github.RepositoryRelease) (string, error) {
	for _, a := range rel.Assets {
		if a.Name != nil && a.ContentType != nil {
			if *a.ContentType == "application/zip" {
				return *a.Name, nil
			}
		}
	}
	return "", errors.New("Asset not found")
}

var Service = new(relver)
