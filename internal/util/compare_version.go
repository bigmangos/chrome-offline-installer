package util

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

func IsNewVersion(new, old string) bool {
	av := parseVersion(new)
	bv := parseVersion(old)
	if len(av) != len(bv) {
		return false
	}
	for i := 0; i < len(av); i++ {
		if av[i] > bv[i] {
			return true
		} else if av[i] < bv[i] {
			return false
		}
	}
	return false
}

func parseVersion(in string) []int {
	r := strings.Split(in, ".")
	v := make([]int, 0, len(r))
	for _, s := range r {
		n, err := strconv.Atoi(s)
		if err != nil {
			slog.Error(fmt.Sprintf("atoi err: %v, str: %s", err, s))
			break
		}
		v = append(v, n)
	}
	return v
}
