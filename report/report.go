package report

import (
	"fmt"
	"gestureflame/domain"
	"sort"
	"strings"
)

func RenderRecord(r domain.Record, events []domain.AuditEvent, attachments []domain.Attachment) string {
	var b strings.Builder
	fmt.Fprintf(&b, "batch=%s\ntitle=%s\nstatus=%s\nresult=%s\nrevision=%d\n", r.BatchCode, r.Title, r.Status, r.Result, r.Revision)
	fmt.Fprintf(&b, "events=%d\nattachments=%d\n", len(events), len(attachments))
	for _, e := range events {
		fmt.Fprintf(&b, "event=%s actor=%s action=%s detail=%s\n", e.ID, e.Actor, e.Action, e.Detail)
	}
	for _, a := range attachments {
		fmt.Fprintf(&b, "attachment=%s type=%s digest=%s\n", a.Name, a.MediaType, a.Digest)
	}
	return b.String()
}
func RenderList(records []domain.Record) string {
	sorted := append([]domain.Record(nil), records...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].BatchCode < sorted[j].BatchCode })
	var b strings.Builder
	for _, r := range sorted {
		fmt.Fprintf(&b, "%s | %s | %s | %s\n", r.BatchCode, r.Title, r.Status, r.Result)
	}
	return b.String()
}
func RenderWorkflow(w domain.Workflow) string {
	return fmt.Sprintf("workflow=%s record=%s name=%s stage=%s owner=%s", w.ID, w.RecordID, w.Name, w.Stage, w.Owner)
}
func StatusLabel(s domain.Status) string {
	switch s {
	case domain.StatusDraft:
		return "Draft"
	case domain.StatusReviewed:
		return "Reviewed"
	case domain.StatusConfirmed:
		return "Confirmed"
	case domain.StatusArchived:
		return "Archived"
	case domain.StatusPublished:
		return "Published"
	default:
		return "Unknown"
	}
}
