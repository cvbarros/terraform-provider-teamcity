package teamcity

import "strings"

func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "404")
}
