package locale

type Locale interface {
	RegisterOnLocaleChangeCallback(func(locale string))
	GetLocale() (string, error)
	SetLocale(token, locale string) (error, bool)
	GetAllowedLocales() []string
}
