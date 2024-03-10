package str

import (
	"strings"
)

func StringContains(slices []string, comparizon string) bool {
	for _, a := range slices {
		if a == comparizon {
			return true
		}
	}

	return false
}

func dumpCount(list []string) map[string]int {

	dupFrequency := make(map[string]int)

	for _, item := range list {
		// check if the item/element exist in the dupFrequency map
		_, exist := dupFrequency[item]

		if exist {
			dupFrequency[item] += 1 // increase counter by 1 if already in the map
		} else {
			dupFrequency[item] = 1 // else start counting from 1
		}
	}
	return dupFrequency
}

func StringContainsPrefix(prefixSlices []string, s string) bool {

	for _, a := range prefixSlices {
		if strings.HasPrefix(s, a) {
			return true
		}
	}

	return false
}
