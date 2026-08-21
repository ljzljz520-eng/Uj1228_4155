package service

import (
	"gestureflame/report"
	"testing"
)

func TestWorkflowImportReport(t *testing.T) {
	app := testService(t)
	r, err := app.Import(testDraft())
	if err != nil {
		t.Fatal(err)
	}
	rec, events, at, err := app.ReportData(r.ID)
	if err != nil || rec.ID != r.ID || len(events) == 0 || len(at) != 1 {
		t.Fatal(err, rec, len(events), len(at))
	}
	text := report.RenderRecord(rec, events, at)
	if len(text) < 20 {
		t.Fatal(text)
	}
}
