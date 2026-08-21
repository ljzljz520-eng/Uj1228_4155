package report

import (
	"gestureflame/domain"
	"strings"
)

func MatchBatch(records []domain.Record, code string) []domain.Record {
	out := []domain.Record{}
	for _, r := range records {
		if strings.EqualFold(r.BatchCode, code) {
			out = append(out, r)
		}
	}
	return out
}
func MatchStatus(records []domain.Record, status domain.Status) []domain.Record {
	out := []domain.Record{}
	for _, r := range records {
		if r.Status == status {
			out = append(out, r)
		}
	}
	return out
}
func CountByStatus(records []domain.Record) map[string]int {
	out := map[string]int{}
	for _, r := range records {
		out[string(r.Status)]++
	}
	return out
}
func EmptyReport() string { return "no records\n" }
