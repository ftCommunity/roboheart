package relver

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/blang/semver"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/package/threadmanager"
	"github.com/google/go-github/v31/github"
)

type relver struct {
	release, prerelease *semver.Version
	lock                sync.Mutex
	last                time.Time
	gh                  *github.Client
	logger              service.LoggerFunc
	error               service.ErrorFunc
	tm                  *threadmanager.ThreadManager
}

type ReleaseVersion interface {
}

func (r *relver) Init(services map[string]service.Service, logger service.LoggerFunc, e service.ErrorFunc) error {
	r.logger = logger
	r.error = e
	r.gh = github.NewClient(nil)
	r.tm = threadmanager.NewThreadManager(r.logger, r.error)
	r.tm.Load("update", r.updateThread)
	r.tm.Start("update")
	return nil
}

func (r *relver) Stop() error                                                { return nil }
func (r *relver) Name() string                                               { return "relver" }
func (r *relver) Dependencies() ([]string, []string)                         { return []string{}, []string{} }
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
	for _, rel := range releases {
		if rel.TagName != nil {
			ver, err := semver.Make(strings.Replace(*rel.TagName, "v", "", -1))
			if err != nil {
				return err
			}
			pr := false
			if rel.Prerelease != nil {
				pr = *rel.Prerelease
			}
			if r.release == nil && !pr {
				r.release = &ver
			}
			if r.prerelease == nil && pr {
				r.prerelease = &ver
			}
		}
	}
	return nil
}

var Service = new(relver)
