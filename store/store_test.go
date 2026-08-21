package store

import (
	"gestureflame/domain"
	"path/filepath"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.db")
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	r := domain.Record{ID: "record-1", BatchCode: "ZX122836", Title: "handoff", Status: domain.StatusArchived, Result: "confirmed flame", Revision: 4, CreatedAt: "a", UpdatedAt: "b"}
	if err := db.SaveRecord(r); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveEvent(domain.AuditEvent{ID: "event-1", RecordID: r.ID, Action: "archive", Actor: "op", Detail: "ok", CreatedAt: "a"}); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveWorkflow(domain.Workflow{ID: "workflow-1", RecordID: r.ID, Name: "handoff", Stage: "done", Owner: "op"}); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveAttachment(domain.Attachment{ID: "attachment-1", RecordID: r.ID, Name: "trace", MediaType: "text/plain", Digest: "d"}); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	got, err := db.GetRecord(r.ID)
	if err != nil || got.BatchCode != r.BatchCode {
		t.Fatalf("%v %#v", err, got)
	}
	events, _ := db.ListEvents(r.ID)
	if len(events) != 1 {
		t.Fatal(len(events))
	}
	wf, err := db.GetWorkflow("workflow-1")
	if err != nil || wf.RecordID != r.ID {
		t.Fatal(err, wf)
	}
	at, _ := db.ListAttachments(r.ID)
	if len(at) != 1 {
		t.Fatal(len(at))
	}
}
