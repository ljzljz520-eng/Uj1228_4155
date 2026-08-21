package cli

import (
	"bytes"
	"gestureflame/clock"
	"gestureflame/service"
	"gestureflame/store"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunImportSearch(t *testing.T) {
	db, _ := store.Open(filepath.Join(t.TempDir(), "x.db"))
	defer db.Close()
	app := service.New(db, clock.New())
	var out bytes.Buffer
	data := `{"batch_code":"ZX12","title":"x","result":"ok","actor":"a","attachment_name":"n","attachment_type":"text","attachment_digest":"d"}`
	if err := Run(app, []string{"import", data}, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "record-") {
		t.Fatal(out.String())
	}
	out.Reset()
	if err := Run(app, []string{"search", "ZX12"}, &out); err != nil || !strings.Contains(out.String(), "ZX12") {
		t.Fatal(err, out.String())
	}
}
