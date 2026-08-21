package service

import (
	"gestureflame/domain"
	"testing"
)

func TestWorkflowCreateReviewArchive(t *testing.T) {
	app := testService(t)
	r, err := app.RunCreateReviewConfirmArchive(testDraft(), func(r domain.Record) error { return nil })
	if err != nil {
		t.Fatal(err)
	}
	if r.Status != domain.StatusArchived {
		t.Fatal(r.Status)
	}
	recalled, events, err := app.Recall(r.ID)
	if err != nil || recalled.Result != "confirmed flame" || len(events) < 4 {
		t.Fatal(err, recalled, len(events))
	}
}
