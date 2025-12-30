package form

import "github.com/freitasmatheusrn/social-fit/pkg/rest"

func HasFieldErrors(errs *rest.ApiErr, fieldName string) bool {
	if errs == nil || len(errs.Causes) == 0 {
		return false
	}

	for _, cause := range errs.Causes {
		if cause.Field == fieldName {
			return true
		}
	}

	return false
}

func GetFieldErrors(errs *rest.ApiErr, fieldName string) []rest.Causes {
	if errs == nil || len(errs.Causes) == 0 {
		return nil
	}

	var fieldErrors []rest.Causes
	for _, cause := range errs.Causes {
		if cause.Field == fieldName {
			fieldErrors = append(fieldErrors, cause)
		}
	}

	return fieldErrors
}
