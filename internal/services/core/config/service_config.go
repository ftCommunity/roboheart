package config

type ServiceConfig struct {
	c        *config
	confname string
}

func (sc *ServiceConfig) GetSections(secType string) ([]string, bool) {
	return sc.c.tree.GetSections(sc.confname, secType)
}

func (sc *ServiceConfig) Get(section, option string) ([]string, bool) {
	return sc.c.tree.Get(sc.confname, section, option)
}

func (sc *ServiceConfig) GetLast(section, option string) (string, bool) {
	return sc.c.tree.GetLast(sc.confname, section, option)
}

func (sc *ServiceConfig) GetBool(section, option string) (bool, bool) {
	return sc.c.tree.GetBool(sc.confname, section, option)
}

func (sc *ServiceConfig) GetBoolDefault(section, option string, def bool) bool {
	if v, ok := sc.GetBool(section, option); ok {
		return v
	} else {
		return def
	}
}

func (sc *ServiceConfig) Set(section, option string, values ...string) bool {
	return sc.c.tree.Set(sc.confname, section, option, values...)
}

func (sc *ServiceConfig) Del(section, option string) { sc.c.tree.Del(sc.confname, section, option) }

func (sc *ServiceConfig) AddSection(section, typ string) error {
	return sc.c.tree.AddSection(sc.confname, section, typ)
}

func (sc *ServiceConfig) DelSection(section string) { sc.c.tree.DelSection(sc.confname, section) }

func newServiceConfig(c *config, confname string) *ServiceConfig {
	return &ServiceConfig{c, confname}
}
