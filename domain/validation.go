package domain

import "strings"

func ValidateRecord(r Record) error {
	if strings.TrimSpace(r.ID) == "" || strings.TrimSpace(r.BatchCode) == "" {
		return ErrMissingField
	}
	if strings.TrimSpace(r.Title) == "" || strings.TrimSpace(r.Result) == "" {
		return ErrMissingField
	}
	if r.Revision < 1 {
		return ErrMissingField
	}
	if !validStatus(r.Status) {
		return ErrInvalidTransition
	}
	return nil
}

func validStatus(s Status) bool {
	switch s {
	case StatusDraft, StatusReviewed, StatusConfirmed, StatusArchived, StatusPublished:
		return true
	default:
		return false
	}
}

func ValidateImport(d ImportDraft) error {
	if strings.TrimSpace(d.BatchCode) == "" || strings.TrimSpace(d.Title) == "" {
		return ErrMissingField
	}
	if strings.TrimSpace(d.Result) == "" || strings.TrimSpace(d.Actor) == "" {
		return ErrMissingField
	}
	if strings.TrimSpace(d.AttachmentName) == "" || strings.TrimSpace(d.AttachmentType) == "" || strings.TrimSpace(d.AttachmentDigest) == "" {
		return ErrMissingField
	}
	if len(d.BatchCode) < 4 {
		return ErrMissingField
	}
	return nil
}

func NormalizeBatchCode(code string) string { return strings.ToUpper(strings.TrimSpace(code)) }
func NormalizeTitle(title string) string    { return strings.Join(strings.Fields(title), " ") }
func NormalizeResult(result string) string  { return strings.TrimSpace(result) }
func IsBatchCode(code string) bool          { return len(NormalizeBatchCode(code)) >= 4 }
func SameResult(a, b string) bool           { return NormalizeResult(a) == NormalizeResult(b) }
