package validators

import (
	"regexp"
)

const EMAIL_REGEX = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

func IsValidateEmail(input string) bool {
	re := regexp.MustCompile(EMAIL_REGEX)

	return re.MatchString(input)
}
