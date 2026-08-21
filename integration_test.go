package gestureflame

import (
	"gestureflame/clock"
	"gestureflame/domain"
	"gestureflame/service"
	"gestureflame/store"
	"path/filepath"
	"testing"
)

func TestBusiness36Regression(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	app := service.New(db, clock.New())
	d := domain.ImportDraft{BatchCode: "ZX122836", Title: "handoff", Result: "confirmed flame", Actor: "operator", AttachmentName: "trace", AttachmentType: "text/plain", AttachmentDigest: "digest"}
	first, err := app.RunCreateReviewConfirmArchive(d, func(domain.Record) error { return nil })
	if err != nil {
		t.Fatal(err)
	}
	recalled, events, err := app.Recall(first.ID)
	if err != nil || recalled.Status != domain.StatusArchived || recalled.Result != d.Result || len(events) < 4 {
		t.Fatalf("initial archive recall lost confirmed state: %v %#v %d", err, recalled, len(events))
	}
	defer func() {
		if recover() != nil {
			t.Errorf("archive should preserve independent confirmed state")
		}
	}()
	_, _ = app.RunCreateReviewConfirmArchive(d, nil)
}
