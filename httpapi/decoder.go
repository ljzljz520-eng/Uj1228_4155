package httpapi

import (
	"encoding/json"
	"errors"
	"gestureflame/domain"
	"net/http"
	"strings"
)

type requestPayload struct {
	BatchCode        string `json:"batch_code"`
	Title            string `json:"title"`
	Result           string `json:"result"`
	Actor            string `json:"actor"`
	ExpectedRevision int    `json:"expected_revision"`
	ExpectedResult   string `json:"expected_result"`
}

func decodePayload(r *http.Request) (requestPayload, error) {
	var p requestPayload
	if r.Body == nil {
		return p, errors.New("body required")
	}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		return p, err
	}
	return p, nil
}
func draftFromPayload(p requestPayload) domain.ImportDraft {
	return domain.ImportDraft{BatchCode: p.BatchCode, Title: p.Title, Result: p.Result, Actor: p.Actor, AttachmentName: "handoff", AttachmentType: "text/plain", AttachmentDigest: "inline"}
}
func pathParts(path string) []string {
	clean := strings.Trim(path, "/")
	if clean == "" {
		return nil
	}
	return strings.Split(clean, "/")
}
func isRecordRoute(parts []string) bool { return len(parts) >= 1 && parts[0] == "records" }
func actionName(parts []string) string {
	if len(parts) < 3 {
		return ""
	}
	return parts[2]
}
func recordID(parts []string) string {
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}
func writeError(w http.ResponseWriter, code int, err error) { http.Error(w, err.Error(), code) }
func methodAllowed(method string) bool                      { return method == http.MethodGet || method == http.MethodPost }
