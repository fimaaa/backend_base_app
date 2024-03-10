package str

import "regexp"

// Sanitize ...
func Sanitize(input string) string {
	reLeadClose := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
	reInside := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	final := reLeadClose.ReplaceAllString(input, "")
	final = reInside.ReplaceAllString(final, " ")

	return final
}
