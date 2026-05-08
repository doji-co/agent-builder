package model

import "strings"

func normalizeName(name string) string {
	var result strings.Builder
	runes := []rune(strings.ReplaceAll(name, "-", "_"))

	for i, r := range runes {
		if r == '_' || r == ' ' {
			if result.Len() > 0 {
				current := result.String()
				if current[len(current)-1] != '_' {
					result.WriteRune('_')
				}
			}
			continue
		}

		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				prev := runes[i-1]
				var next rune
				if i+1 < len(runes) {
					next = runes[i+1]
				}
				if (prev >= 'a' && prev <= 'z') || (prev >= '0' && prev <= '9') || (prev >= 'A' && prev <= 'Z' && next >= 'a' && next <= 'z') {
					if result.Len() > 0 {
						current := result.String()
						if current[len(current)-1] != '_' {
							result.WriteRune('_')
						}
					}
				}
			}
			r += 'a' - 'A'
		}

		result.WriteRune(r)
	}

	return strings.Trim(result.String(), "_")
}

func className(name string) string {
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return r == '-' || r == '_' || r == ' '
	})

	var result strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}

		var wroteFirst bool
		for _, r := range part {
			if !wroteFirst {
				if r >= 'a' && r <= 'z' {
					r -= 'a' - 'A'
				}
				result.WriteRune(r)
				wroteFirst = true
				continue
			}

			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
				result.WriteRune(r)
			}
		}
	}

	return result.String()
}
