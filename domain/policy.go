package domain

import (
	"fmt"
	"strings"
)

type FieldIssue struct {
	Field   string
	Message string
}
type ReviewChecklist struct {
	Identity  bool
	Evidence  bool
	Result    bool
	Ownership bool
}

func (c ReviewChecklist) Complete() bool { return c.Identity && c.Evidence && c.Result && c.Ownership }
func (c ReviewChecklist) Missing() []string {
	out := []string{}
	if !c.Identity {
		out = append(out, "identity")
	}
	if !c.Evidence {
		out = append(out, "evidence")
	}
	if !c.Result {
		out = append(out, "result")
	}
	if !c.Ownership {
		out = append(out, "ownership")
	}
	return out
}
func ValidateFields(r Record) []FieldIssue {
	issues := []FieldIssue{}
	if strings.TrimSpace(r.ID) == "" {
		issues = append(issues, FieldIssue{"id", "missing id"})
	}
	if !IsBatchCode(r.BatchCode) {
		issues = append(issues, FieldIssue{"batch_code", "invalid batch code"})
	}
	if strings.TrimSpace(r.Title) == "" {
		issues = append(issues, FieldIssue{"title", "missing title"})
	}
	if strings.TrimSpace(r.Result) == "" {
		issues = append(issues, FieldIssue{"result", "missing result"})
	}
	if r.Revision < 1 {
		issues = append(issues, FieldIssue{"revision", "revision must be positive"})
	}
	return issues
}
func ExplainIssues(issues []FieldIssue) string {
	if len(issues) == 0 {
		return "valid"
	}
	parts := make([]string, 0, len(issues))
	for _, issue := range issues {
		parts = append(parts, fmt.Sprintf("%s: %s", issue.Field, issue.Message))
	}
	return strings.Join(parts, "; ")
}
func ValidateActor(actor string) bool {
	actor = strings.TrimSpace(actor)
	return len(actor) >= 2 && !strings.ContainsAny(actor, "\n\r")
}
func ValidateAttachment(a Attachment) []FieldIssue {
	out := []FieldIssue{}
	if strings.TrimSpace(a.RecordID) == "" {
		out = append(out, FieldIssue{"record_id", "missing record"})
	}
	if strings.TrimSpace(a.Name) == "" {
		out = append(out, FieldIssue{"name", "missing name"})
	}
	if strings.TrimSpace(a.MediaType) == "" {
		out = append(out, FieldIssue{"media_type", "missing type"})
	}
	if strings.TrimSpace(a.Digest) == "" {
		out = append(out, FieldIssue{"digest", "missing digest"})
	}
	return out
}
func ValidateWorkflow(w Workflow) []FieldIssue {
	out := []FieldIssue{}
	if w.ID == "" {
		out = append(out, FieldIssue{"id", "missing id"})
	}
	if w.RecordID == "" {
		out = append(out, FieldIssue{"record_id", "missing record"})
	}
	if w.Name == "" {
		out = append(out, FieldIssue{"name", "missing name"})
	}
	if w.Stage == "" {
		out = append(out, FieldIssue{"stage", "missing stage"})
	}
	return out
}
func IsAllowedAction(status Status, action string) bool {
	switch action {
	case "review":
		return status == StatusDraft
	case "confirm":
		return status == StatusReviewed
	case "archive":
		return status == StatusConfirmed
	case "publish":
		return status == StatusArchived
	default:
		return false
	}
}
func ActionForStatus(status Status) string {
	switch status {
	case StatusDraft:
		return "review"
	case StatusReviewed:
		return "confirm"
	case StatusConfirmed:
		return "archive"
	case StatusArchived:
		return "publish"
	default:
		return "none"
	}
}
func LifecycleSummary(r Record) string {
	return fmt.Sprintf("%s is %s at revision %d", r.BatchCode, r.Status, r.Revision)
}
func CloneEvents(events []AuditEvent) []AuditEvent {
	out := make([]AuditEvent, len(events))
	copy(out, events)
	return out
}
