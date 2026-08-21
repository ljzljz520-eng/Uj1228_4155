package httpapi

import (
	"gestureflame/clock"
	"gestureflame/domain"
	"gestureflame/service"
	"gestureflame/store"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestHTTPCreateAndGet(t *testing.T) {
	db, _ := store.Open(filepath.Join(t.TempDir(), "x.db"))
	defer db.Close()
	app := service.New(db, clock.New())
	h := WithHealth(Routes(New(app)))
	body := strings.NewReader(`{"batch_code":"ZX12","title":"x","result":"ok","actor":"a"}`)
	req := httptest.NewRequest(http.MethodPost, "/records", body)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatal(rec.Code, rec.Body.String())
	}
	req = httptest.NewRequest(http.MethodGet, "/records?q=ZX12", nil)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "ZX12") {
		t.Fatal(rec.Code, rec.Body.String())
	}
	_ = io.Discard
	_ = domain.StatusDraft
}
