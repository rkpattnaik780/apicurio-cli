package contextcmdutil

import (
	"regexp"

	"github.com/apicurio/apicurio-cli/pkg/core/errors"
	"github.com/apicurio/apicurio-cli/pkg/core/localize"
	"github.com/apicurio/apicurio-cli/pkg/core/servicecontext"
	"github.com/apicurio/apicurio-cli/pkg/shared/contextutil"
)

// Validator is a type for validating service context inputs
type Validator struct {
	Localizer  localize.Localizer
	SvcContext *servicecontext.Context
}

const (
	legalNameChars = "^[a-z]([-a-z0-9]*[a-z0-9])?$"
)

// ValidateName validates the name of the context
func (v *Validator) ValidateName(val interface{}) error {

	name, ok := val.(string)
	if !ok {
		return errors.NewCastError(val, "string")
	}

	if len(name) < 1 {
		return v.Localizer.MustLocalizeError("context.common.validation.name.error.required")
	}

	matched, _ := regexp.Match(legalNameChars, []byte(name))

	if matched {
		return nil
	}

	return v.Localizer.MustLocalizeError("context.common.validation.name.error.invalidChars", localize.NewEntry("Name", name))
}

// ValidateNameIsAvailable validates if the name provided is a unique context name
func (v *Validator) ValidateNameIsAvailable(name string) error {

	context, _ := contextutil.GetContext(v.SvcContext, v.Localizer, name)
	if context != nil {
		return v.Localizer.MustLocalizeError("context.create.log.alreadyExists", localize.NewEntry("Name", name))
	}

	return nil
}
