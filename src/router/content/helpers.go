package content

import "strconv"

func parseIntQueryParam(param string) (int, error) {
	if len(param) != 0 {
		return strconv.Atoi(param)
	}
	return 0, nil
}
