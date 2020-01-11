package discord

func startsWith(source, pattern string) bool {
	if len(source) < len(pattern) || len(source) < 1 {
		return false
	}
	if source[:len(pattern)] != pattern {
		return false
	}
	return true
}
