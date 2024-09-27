package utils

import "errors"

func UnwrapRecursive(err error, targetError []error) error {
	for {
		for _, target := range targetError {
			if errors.Is(err, target) {
				return err
			}
		}
		if e, ok := err.(interface{ Unwrap() error }); ok {
			err = e.Unwrap()
		} else {
			break
		}
	}
	return err
}
