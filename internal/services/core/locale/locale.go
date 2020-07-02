package locale

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/services/core/acm"
	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/api"
	fileperm "github.com/ftCommunity/roboheart/package/filepermissions"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
	"github.com/gorilla/mux"
	"github.com/thoas/go-funk"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/ftCommunity/roboheart/internal/service"
)

const (
	PERMISSION = "locale"
	LOCALEPATH = "/etc/locale"
)

var (
	LOCALES = []string{"en_US", "de_DE", "fr_FR", "nl_NL"}
)

type locale struct {
	logger    service.LoggerFunc
	error     service.ErrorFunc
	locales   []string
	callbacks []func(string)
	lock      sync.Mutex
	web       web.Web
	mux       *mux.Router
	acm       acm.ACM
}

type Locale interface {
	RegisterOnLocaleChangeCallback(func(locale string))
	GetLocale() (string, error)
	SetLocale(token, locale string) (error, bool)
	GetAllowedLocales() []string
}

func (l *locale) Init(services map[string]service.Service, logger service.LoggerFunc, e service.ErrorFunc) error {
	l.logger = logger
	l.error = e
	if err := servicehelpers.CheckMainDependencies(l, services); err != nil {
		return err
	}
	var ok bool
	l.acm, ok = services["acm"].(acm.ACM)
	if !ok {
		return errors.New("Type assertion error")
	}
	l.acm.RegisterPermission(PERMISSION, map[string]bool{"user": true, "app": false}, map[string]string{})
	for _, ln := range LOCALES {
		if !strings.Contains(ln, ".") {
			ln = ln + ".UTF-8"
		}
		l.locales = append(l.locales, ln)
	}
	return nil
}

func (l *locale) Stop() error {
	l.lock.Lock()
	return nil
}

func (l *locale) Name() string                       { return "locale" }
func (l *locale) Dependencies() ([]string, []string) { return []string{"acm"}, []string{"web"} }

func (l *locale) SetAdditionalDependencies(services map[string]service.Service) error {
	if err := servicehelpers.CheckAdditionalDependencies(l, services); err != nil {
		return err
	}
	var ok bool
	l.web, ok = services["web"].(web.Web)
	if !ok {
		return errors.New("Type assertion error")
	}
	l.configureWeb()
	return nil
}

func (l *locale) configureWeb() {
	l.mux = l.web.RegisterServiceAPI(l)
	l.mux.HandleFunc("/locale", func(w http.ResponseWriter, r *http.Request) {
		if locale, err := l.GetLocale(); err != nil {
			api.ErrorResponseWriter(w, 500, err)
		} else {
			api.ResponseWriter(w, locale)
		}
	}).Methods("GET")
	l.mux.HandleFunc("/locale", func(w http.ResponseWriter, r *http.Request) {
		data := &struct {
			api.TokenRequest
			Locale string `json:"locale"`
		}{}
		if !api.RequestLoader(r, w, data) {
			return
		}
		if err, uae := l.SetLocale(data.Token, data.Locale); err != nil {
			code := 500
			if uae {
				code = 403
			}
			api.ErrorResponseWriter(w, code, err)
		} else {
			api.ResponseWriter(w, nil)
		}
	}).Methods("POST")
	l.mux.HandleFunc("/allowed", func(w http.ResponseWriter, r *http.Request) {
		api.ResponseWriter(w, l.GetAllowedLocales())
	}).Methods("GET")
}

func (l *locale) RegisterOnLocaleChangeCallback(cb func(locale string)) {
	l.callbacks = append(l.callbacks, cb)
}

func (l *locale) SetLocale(token, locale string) (error, bool) {
	if err, uae := l.acm.CheckTokenPermission(token, PERMISSION); err != nil {
		return err, uae
	}
	if !funk.ContainsString(l.locales, locale) {
		return errors.New("Locale unknown"), false
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	raw, err := ioutil.ReadFile(LOCALEPATH)
	if err != nil {
		return err, false
	}
	olddata := strings.Split(string(raw), "\n")
	data := []string{}
	for _, line := range olddata {
		if !strings.HasPrefix(line, "LC_ALL") {
			data = append(data, line)
		}
	}
	data = append(data, "LC_ALL=\""+locale+"\"")
	if err := ioutil.WriteFile(LOCALEPATH, []byte(strings.Join(data, "\n")), fileperm.OS_U_RW_G_RW_O_R); err != nil {
		return err, false
	}
	go l.runCallbacks(locale)
	return nil, false
}

func (l *locale) GetLocale() (string, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	raw, err := ioutil.ReadFile(LOCALEPATH)
	if err != nil {
		return "", err
	}
	locale := ""
	data := strings.Split(string(raw), "\n")
	for _, line := range data {
		if strings.HasPrefix(line, "LC_ALL") {
			locale = strings.Replace(line, "LC_ALL=", "", 1)
		}
	}
	locale = strings.Replace(locale, "\"", "", 2)
	if locale == "" {
		return "", errors.New("No locale set")
	}
	return locale, nil
}

func (l *locale) GetAllowedLocales() []string {
	return l.locales
}

func (l *locale) runCallbacks(ln string) {
	for _, c := range l.callbacks {
		go c(ln)
	}
}

var Service = new(locale)
