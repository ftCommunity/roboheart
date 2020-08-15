package pkgmanager

import (
	"encoding/json"
	"errors"
	"github.com/blang/semver"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/acm"
	"github.com/ftCommunity/roboheart/internal/services/core/config"
	"github.com/ftCommunity/roboheart/internal/services/core/deviceinfo"
	"github.com/ftCommunity/roboheart/internal/services/core/fwver"
	"github.com/ftCommunity/roboheart/internal/services/core/web"
	fileperm "github.com/ftCommunity/roboheart/package/filepermissions"
	"github.com/ftCommunity/roboheart/package/threadmanager"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

const (
	PERMISSION_BASE         = "pkgmanager"
	PERMISSION_INSTALL      = PERMISSION_BASE + "." + "install"
	PERMISSION_REMOVE       = PERMISSION_BASE + "." + "remove"
	PERMISSION_UPDATE       = PERMISSION_BASE + "." + "update"
	PERMISSION_GETAVAILABLE = PERMISSION_BASE + "." + "getavailable"
	CONFIG_SECTION          = ""
	CONFIG_TYPE             = "pkgmanager"
	PATH_BASE               = "/opt/ftc"
	PATH_PKG                = PATH_BASE + "/" + "packages"
	PATH_DATA               = PATH_BASE + "/" + "data"
	MANIFEST_NAME           = "manifest.json"
)

type pkgmanager struct {
	logger           service.LoggerFunc
	error            service.ErrorFunc
	tm               *threadmanager.ThreadManager
	acm              acm.ACM
	config           config.Config
	sconfig          *config.ServiceConfig
	deviceinfo       deviceinfo.DeviceInfo
	platform, device string
	fwver            fwver.FWVer
	firmware         semver.Version
	web              web.Web
	packages         map[string]extendedPackage
	treelock         sync.Mutex
}

type PkgManager interface {
}

func (p *pkgmanager) Init(services map[string]service.Service, logger service.LoggerFunc, e service.ErrorFunc) {
	p.logger = logger
	p.error = e
	var ok bool
	p.acm, ok = services["acm"].(acm.ACM)
	if !ok {
		e(errors.New("Type assertion error"))
	}
	p.acm.RegisterPermission(PERMISSION_INSTALL, map[string]bool{"user": true, "app": false}, map[string]string{})
	p.acm.RegisterPermission(PERMISSION_REMOVE, map[string]bool{"user": true, "app": false}, map[string]string{})
	p.acm.RegisterPermission(PERMISSION_UPDATE, map[string]bool{"user": true, "app": false}, map[string]string{})
	p.acm.RegisterPermission(PERMISSION_GETAVAILABLE, map[string]bool{"user": true, "app": true}, map[string]string{})
	p.config, ok = services["config"].(config.Config)
	if !ok {
		e(errors.New("Type assertion error"))
	}
	p.sconfig = p.config.GetServiceConfig(p)
	if err := p.sconfig.AddSection(CONFIG_SECTION, CONFIG_TYPE); err != nil {
		e(err)
	}
	p.deviceinfo, ok = services["deviceinfo"].(deviceinfo.DeviceInfo)
	if !ok {
		e(errors.New("Type assertion error"))
	}
	p.device, p.platform = p.deviceinfo.GetDevice(), p.deviceinfo.GetPlatform()
	p.fwver, ok = services["fwver"].(fwver.FWVer)
	if !ok {
		e(errors.New("Type assertion error"))
	}
	p.firmware = p.fwver.Get()
	if err := os.MkdirAll(PATH_PKG, fileperm.OS_U_RW_G_RW_O_R); err != nil {
		e(err)
	}
	if err := os.MkdirAll(PATH_DATA, fileperm.OS_U_RW_G_RW_O_R); err != nil {
		e(err)
	}
	p.packages = make(map[string]extendedPackage)
	p.tm = threadmanager.NewThreadManager(p.logger, p.error)
	go p.reloadAll()
}

func (p *pkgmanager) Stop() {
	p.tm.StopAll()
}

func (p *pkgmanager) Name() string {
	return "pkgmanager"
}

func (p *pkgmanager) Dependencies() service.ServiceDependencies {
	return service.ServiceDependencies{Deps: []string{"acm", "config", "deviceinfo", "fwver"}, ADeps: []string{"web"}}
}

func (p *pkgmanager) SetAdditionalDependencies(services map[string]service.Service) {
	p.web = services["web"].(web.Web)
	p.configureWeb()
}

func (p *pkgmanager) UnsetAdditionalDependencies([]string) {}

func (p *pkgmanager) configureWeb() {
}

func (p *pkgmanager) loadPackageManifest(pkg, variant string) error {
	manifestpath := path.Join(PATH_PKG, pkg, variant, MANIFEST_NAME)
	if _, err := os.Stat(manifestpath); os.IsNotExist(err) {
		return err
	}
	file, err := ioutil.ReadFile("test.json")
	if err != nil {
		return err
	}
	data := extendedVariant{}
	if err = json.Unmarshal([]byte(file), &data); err != nil {
		return err
	}
	if _, ok := p.packages[pkg]; !ok {
		epkg := extendedPackage{}
		epkg.Id = pkg
		p.packages[pkg] = epkg
	}
	p.packages[pkg].Variants[variant] = data
	return nil
}

func (p *pkgmanager) loadManifests() error {
	pkgs, err := ioutil.ReadDir(PATH_PKG)
	if err != nil {
		return err
	}
pkgloop:
	for _, pd := range pkgs {
		if !pd.IsDir() {
			continue pkgloop
		}
		pkgname := p.Name()
		variants, err := ioutil.ReadDir(path.Join(PATH_PKG, pkgname))
		if err != nil {
			return err
		}
	variantloop:
		for _, vd := range variants {
			if !vd.IsDir() {
				continue variantloop
			}
			variantname := vd.Name()
			p.loadPackageManifest(pkgname, variantname)
		}
	}
	return nil
}

func (p *pkgmanager) reloadAll() error {
	return nil
}

var Service = new(pkgmanager)
