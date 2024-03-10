package str

import (
	"regexp"
	"strings"
)

func AddZeroCharToPhone(phone string) string {
	isMatch, _ := regexp.MatchString(`^[0{1}]`, phone)
	if !isMatch {
		isMatchPlus, _ := regexp.MatchString(`^\+[0-9]{2}`, phone)
		if isMatchPlus {
			re := regexp.MustCompile(`^\+[0-9]{2}`)
			s := re.ReplaceAllString(phone, `0`)

			return s
		}

		return "0" + phone
	}

	return phone
}

func PhoneConvertToAbbv(phone string) string {
	isMatch, _ := regexp.MatchString(`^[0{1}]`, phone)
	if isMatch {
		re := regexp.MustCompile(`^[0{1}]`)
		s := re.ReplaceAllString(phone, `+62`)

		return s
	}

	return phone
}

func PhoneConvertToAbbvWithoutPlus(phone string) string {
	isMatch, _ := regexp.MatchString(`^[0{1}]`, phone)
	if isMatch {
		re := regexp.MustCompile(`^[0{1}]`)
		s := re.ReplaceAllString(phone, "62")

		return s
	}

	return phone
}

func Replacer(source string, replacer *strings.Replacer) string {
	return replacer.Replace(source)
}
