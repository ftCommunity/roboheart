package acm

type ACM interface {
	RegisterPermission(name string, defaults map[string]bool, desc map[string]string) error
	CreateToken(defaults []string, layers ...map[string]bool) (string, error)
	UpdateToken(id string, layers ...map[string]bool) error
	GetToken(id string) (*token, error)
	CheckTokenPermission(token string, permission string) (error, bool)
	GetPermissionDescription(name string) (map[string]string, error)
}
