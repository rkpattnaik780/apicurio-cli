package flagutil

import (
	"fmt"
	"sort"

	"github.com/apicurio/apicurio-cli/pkg/core/localize"
)

// IsValidInput checks if the input value is in the range of valid values
func IsValidInput(input string, validValues ...string) bool {
	for _, b := range validValues {
		if input == b {
			return true
		}
	}

	return false
}

// FlagDescription creates a flag description and adds a list of valid options (if any)
func FlagDescription(localizer localize.Localizer, messageID string, validOptions ...string) string {
	// ensure consistent order
	sort.Strings(validOptions)

	description := localizer.MustLocalize(messageID)

	var chooseFrom string
	if len(validOptions) > 0 {
		if description[len(description)-1:] != "." {
			description += "."
		}
		chooseFrom = localizer.MustLocalize("flag.common.chooseFrom")

		for i, val := range validOptions {
			chooseFrom += fmt.Sprintf("\"%v\"", val)
			if i < len(validOptions)-1 {
				chooseFrom += ", "
			}
		}
	}

	return fmt.Sprintf("%v %v", description, chooseFrom)
}

// StringInSlice checks if a string is in a slice.
// It is used when array of strings is passed as flag value or argument
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
