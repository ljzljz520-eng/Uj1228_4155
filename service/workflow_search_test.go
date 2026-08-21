package service

import (
	"gestureflame/domain"
	"testing"
)

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	app := testService(t)
	r, err := app.RunCreateReviewConfirmArchive(testDraft(), func(r domain.Record) error { return nil })
	if err != nil {
		t.Fatal(err)
	}
	changed, err := app.Collaborate(r.ID, "partner", "independent result", r.Revision)
	if err != nil {
		t.Fatal(err)
	}
	if changed.Result != "independent result" {
		t.Fatal(changed.Result)
	}
	if _, err := app.PublishWithConfirmation(r.ID, "partner", "independent result"); err != nil {
		t.Fatal(err)
	}
}
