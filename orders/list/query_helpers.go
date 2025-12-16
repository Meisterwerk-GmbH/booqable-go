package list

import (
	"net/url"
	"strings"
)

func ensureInclude(v url.Values, value string) {
	if v == nil {
		return
	}
	current := v.Get("include")
	if current == "" {
		v.Set("include", value)
		return
	}
	parts := strings.Split(current, ",")
	for _, p := range parts {
		if strings.TrimSpace(p) == value {
			return
		}
	}
	v.Set("include", current+","+value)
}

func cloneValues(v url.Values) url.Values {
	if v == nil {
		return url.Values{}
	}
	out := make(url.Values, len(v))
	for key, vals := range v {
		cp := make([]string, len(vals))
		copy(cp, vals)
		out[key] = cp
	}
	return out
}
