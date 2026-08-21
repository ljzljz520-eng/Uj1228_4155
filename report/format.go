package report

import (
	"fmt"
	"gestureflame/domain"
	"strings"
)

func RenderTimeline(events []domain.AuditEvent) string {
	var b strings.Builder
	for i, e := range events {
		fmt.Fprintf(&b, "%02d %s %s by %s: %s\n", i+1, e.CreatedAt, e.Action, e.Actor, e.Detail)
	}
	return b.String()
}
func RenderChecklist(check domain.ReviewChecklist) string {
	if check.Complete() {
		return "checklist=complete"
	}
	return "checklist=missing:" + strings.Join(check.Missing(), ",")
}
func RenderMetrics(values map[string]int) string {
	keys := []string{"records", "events", "workflows", "attachments"}
	parts := []string{}
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%d", key, values[key]))
	}
	return strings.Join(parts, " ")
}
func RenderIssues(issues []domain.FieldIssue) string {
	if len(issues) == 0 {
		return "no issues"
	}
	parts := []string{}
	for _, issue := range issues {
		parts = append(parts, issue.Field+":"+issue.Message)
	}
	return strings.Join(parts, " | ")
}
func RenderSearch(records []domain.Record, query string) string {
	if len(records) == 0 {
		return "no matches for " + query
	}
	return fmt.Sprintf("%d matches for %s\n%s", len(records), query, RenderList(records))
}
func RenderArchiveNote(r domain.Record) string {
	return fmt.Sprintf("Archive %s at revision %d with result %q", r.BatchCode, r.Revision, r.Result)
}
func JoinReports(reports ...string) string { return strings.Join(reports, "\n---\n") }
