package config

type serviceConfig struct {
	c        *config
	confname string
}

func (sc *serviceConfig) GetSections(secType string) ([]string, bool) {
	return sc.c.tree.GetSections(sc.confname, secType)
}

func (sc *serviceConfig) Get(section, option string) ([]string, bool) {
	return sc.c.tree.Get(sc.confname, section, option)
}

func (sc *serviceConfig) GetLast(section, option string) (string, bool) {
	return sc.c.tree.GetLast(sc.confname, section, option)
}

func (sc *serviceConfig) GetBool(section, option string) (bool, bool) {
	return sc.c.tree.GetBool(sc.confname, section, option)
}

func (sc *serviceConfig) Set(section, option string, values ...string) bool {
	return sc.c.tree.Set(sc.confname, section, option, values...)
}

func (sc *serviceConfig) Del(section, option string) { sc.c.tree.Del(sc.confname, section, option) }

func (sc *serviceConfig) AddSection(section, typ string) error {
	return sc.c.tree.AddSection(sc.confname, section, typ)
}

func (sc *serviceConfig) DelSection(section string) { sc.c.tree.DelSection(sc.confname, section) }

func newServiceConfig(c *config, confname string) *serviceConfig {
	return &serviceConfig{c, confname}
}
