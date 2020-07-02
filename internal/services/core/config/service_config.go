package config

import "strconv"

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

func (sc *ServiceConfig) GetInt(section, option string) (int, bool, error) {
	if raw, ok := sc.GetLast(section, option); ok {
		if v, err := strconv.Atoi(raw); err == nil {
			return v, true, nil
		} else {
			return 0, true, err
		}
	} else {
		return 0, false, nil
	}
}

func (sc *ServiceConfig) GetIntDefault(section, option string, def int) (int, error) {
	if v, ok, err := sc.GetInt(section, option); err == nil {
		if ok {
			return v, nil
		} else {
			return def, nil
		}
	} else {
		return 0, err
	}
}

func (sc *ServiceConfig) GetStringDefault(section, option string, def string) string {
	if v, ok := sc.GetLast(section, option); ok {
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
