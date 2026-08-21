package report

import (
	"gestureflame/domain"
	"strings"
	"testing"
)

func TestRenderRecord(t *testing.T) {
	text := RenderRecord(domain.Record{BatchCode: "ZX", Title: "t", Status: domain.StatusArchived, Result: "ok", Revision: 2}, []domain.AuditEvent{{ID: "e", Action: "archive", Actor: "a", Detail: "d"}}, []domain.Attachment{{Name: "x", MediaType: "text"}})
	if !strings.Contains(text, "status=archived") || !strings.Contains(text, "event=e") {
		t.Fatal(text)
	}
}
func TestFilters(t *testing.T) {
	rs := []domain.Record{{BatchCode: "A", Status: domain.StatusDraft}, {BatchCode: "B", Status: domain.StatusArchived}}
	if len(MatchStatus(rs, domain.StatusArchived)) != 1 {
		t.Fatal()
	}
	if CountByStatus(rs)["draft"] != 1 {
		t.Fatal()
	}
}
