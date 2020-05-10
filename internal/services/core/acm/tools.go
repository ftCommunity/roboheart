package acm

import "errors"

var (
	NotPermittedError = errors.New("Operation not permitted")
)

func CheckTokenPermission(acm ACM, token string, permission string) error {
	t, err := acm.GetToken(token)
	if err != nil {
		return err
	}
	perm, err := t.GetPermission(permission)
	if err != nil {
		return err
	}
	if !perm {
		return NotPermittedError
	}
	return nil
}

func IsUserError(err error) bool {
	return err == NotPermittedError || err == TokenNotFoundError
}
