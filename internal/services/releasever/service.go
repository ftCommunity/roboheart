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
	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
	"github.com/ftCommunity/roboheart/package/threadmanager"
	"github.com/google/go-github/v31/github"
	"github.com/gorilla/mux"
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
	web                 web.Web
	mux                 *mux.Router
}

func (r *relver) Init(services map[string]service.Service, logger service.LoggerFunc, e service.ErrorFunc) {
	r.logger = logger
	r.error = e
	if err := servicehelpers.CheckMainDependencies(r, services); err != nil {
		e(err)
	}
	if err := servicehelpers.InitializeDependencies(services, servicehelpers.ServiceInitializers{"acm": r.initSvcAcm}); err != nil {
		e(err)
	}
	r.gh = github.NewClient(nil)
	r.tm = threadmanager.NewThreadManager(r.logger, r.error)
	r.tm.Load("update", r.updateThread)
	r.tm.Start("update")
}

func (r *relver) Stop() {
	r.tm.StopAll()
}

func (r *relver) Name() string { return "relver" }
func (r *relver) Dependencies() service.ServiceDependencies {
	return service.ServiceDependencies{Deps: []string{"acm"}, ADeps: []string{"web"}}
}
func (r *relver) SetAdditionalDependencies(services map[string]service.Service) {
	servicehelpers.InitializeAdditionalDependencies(services, servicehelpers.AdditionalServiceInitializers{"web": r.initSvcWeb})

}
func (r *relver) UnsetAdditionalDependencies([]string) {}

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

func (r *relver) Update(token string) (error, bool) {
	if err, uae := r.acm.CheckTokenPermission(token, PERMISSION_UPDATE); err != nil {
		return err, uae
	}
	return r.getReleaseData(), false
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
	if r.Download, r.Size, err = getAssetURL(rel); err != nil {
		return nil, err
	}
	return r, nil
}

func getAssetURL(rel *github.RepositoryRelease) (string, int, error) {
	for _, a := range rel.Assets {
		if a.Name != nil && a.ContentType != nil && a.Size != nil {
			if *a.ContentType == "application/zip" {
				return *a.Name, *a.Size, nil
			}
		}
	}
	return "", 0, errors.New("Asset not found")
}
