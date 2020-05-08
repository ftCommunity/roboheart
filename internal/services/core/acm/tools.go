package acm

import "errors"

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
		return errors.New("Operation not permitted")
	}
	return nil
}
