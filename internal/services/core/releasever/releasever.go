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

type relver struct {
	release, prerelease *release
	releases            []release
	lock                sync.Mutex
	last                time.Time
	gh                  *github.Client
	logger              service.LoggerFunc
	error               service.ErrorFunc
	tm                  *threadmanager.ThreadManager
	acm                 acm.ACM
}

type ReleaseVersion interface {

type release struct {
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
	r.releases = make([]release, 0)
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

func newRelease(rel *github.RepositoryRelease) (*release, error) {
	r := new(release)
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
