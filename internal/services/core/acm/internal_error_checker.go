package acm

func isUserError(err error) bool {
	return err == NotPermittedError || err == TokenNotFoundError
}
