package domain

import "testing"

func TestValidateAndTransition(t *testing.T) {
	r := Record{ID: "r1", BatchCode: "ZX1", Title: "array", Result: "stable", Status: StatusDraft, Revision: 1}
	if err := ValidateRecord(r); err != nil {
		t.Fatal(err)
	}
	for _, next := range []Status{StatusReviewed, StatusConfirmed, StatusArchived, StatusPublished} {
		if err := Transition(&r, next); err != nil {
			t.Fatal(err)
		}
	}
	if !r.IsTerminal() {
		t.Fatal(r.Status)
	}
}
func TestImportValidation(t *testing.T) {
	if err := ValidateImport(ImportDraft{BatchCode: "ZX12", Title: "x", Result: "ok", Actor: "a", AttachmentName: "n", AttachmentType: "text", AttachmentDigest: "d"}); err != nil {
		t.Fatal(err)
	}
}
