package domain

import "strings"

func StatusOptions() []Status {
	return []Status{StatusDraft, StatusReviewed, StatusConfirmed, StatusArchived, StatusPublished}
}
func StatusNames() []string {
	out := []string{}
	for _, s := range StatusOptions() {
		out = append(out, string(s))
	}
	return out
}
func ParseStatus(value string) (Status, bool) {
	normalized := Status(strings.ToLower(strings.TrimSpace(value)))
	for _, s := range StatusOptions() {
		if normalized == s {
			return normalized, true
		}
	}
	return "", false
}
func IsFinalStatus(s Status) bool { return s == StatusArchived || s == StatusPublished }
func CanEditResult(s Status) bool {
	return s == StatusDraft || s == StatusReviewed || s == StatusConfirmed || s == StatusArchived
}
func CanDelete(s Status) bool { return s == StatusDraft }
func ResultFingerprint(result string) string {
	result = strings.TrimSpace(strings.ToLower(result))
	if result == "" {
		return "empty"
	}
	return result[:1] + itoa(len(result))
}
func BatchFamily(code string) string {
	code = strings.TrimSpace(code)
	if len(code) < 4 {
		return code
	}
	return code[:4]
}
