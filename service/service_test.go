package service

import (
	"gestureflame/clock"
	"gestureflame/domain"
	"gestureflame/store"
	"path/filepath"
	"testing"
)

func testService(t *testing.T) *Service {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return New(db, clock.New())
}
func testDraft() domain.ImportDraft {
	return domain.ImportDraft{BatchCode: "ZX122836", Title: "handoff", Result: "confirmed flame", Actor: "operator", AttachmentName: "trace", AttachmentType: "text/plain", AttachmentDigest: "digest"}
}
