package helper

import "strings"

func GenerateBasePagingMap(searchParam string, limit, page int, sortBy, sortType string) map[string]interface{}{

	dataResult := map[string]interface{}{
		"search_param" 	: strings.TrimSpace(searchParam),
		"limit"			: limit,
		"page" 			: page,
		"sort_by"		: strings.TrimSpace(sortBy),
		"sort_type"		: strings.TrimSpace(sortType),
	}

	return dataResult
}
