package authentication

type Options func(*Authentication)

func WithValidator(validator Validator) Options {
	return func(auth *Authentication) {
		auth.validator = validator
	}
}

func WithIgnorePaths(paths []string) Options {
	return func(auth *Authentication) {
		auth.ignorePaths = paths
	}
}
