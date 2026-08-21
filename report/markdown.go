package report

import (
	"gestureflame/domain"
	"strings"
)

func MarkdownRecord(r domain.Record) string {
	var b strings.Builder
	b.WriteString("## " + r.BatchCode + "\n\n")
	b.WriteString("- Title: " + r.Title + "\n")
	b.WriteString("- Status: " + string(r.Status) + "\n")
	b.WriteString("- Result: " + r.Result + "\n")
	b.WriteString("- Revision: " + fmtRevision(r.Revision) + "\n")
	return b.String()
}
func fmtRevision(value int) string {
	if value < 0 {
		return "invalid"
	}
	digits := ""
	if value == 0 {
		return "0"
	}
	for value > 0 {
		digits = string(rune('0'+value%10)) + digits
		value /= 10
	}
	return digits
}
func MarkdownEvents(events []domain.AuditEvent) string {
	var b strings.Builder
	b.WriteString("| Action | Actor | Detail |\n|---|---|---|\n")
	for _, e := range events {
		b.WriteString("| " + e.Action + " | " + e.Actor + " | " + e.Detail + " |\n")
	}
	return b.String()
}
func MarkdownAttachments(items []domain.Attachment) string {
	var b strings.Builder
	for _, a := range items {
		b.WriteString("- " + a.Name + " (" + a.MediaType + ") " + a.Digest + "\n")
	}
	return b.String()
}
func PlainLines(values []string) string { return strings.Join(values, "\n") + "\n" }
