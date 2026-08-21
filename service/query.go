package service

import (
	"gestureflame/domain"
	"sort"
)

func SortByRevision(records []domain.Record) []domain.Record {
	out := append([]domain.Record(nil), records...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Revision == out[j].Revision {
			return out[i].ID < out[j].ID
		}
		return out[i].Revision > out[j].Revision
	})
	return out
}
func FilterStatus(records []domain.Record, status domain.Status) []domain.Record {
	out := []domain.Record{}
	for _, r := range records {
		if r.Status == status {
			out = append(out, r)
		}
	}
	return out
}
func Latest(records []domain.Record) (domain.Record, bool) {
	if len(records) == 0 {
		return domain.Record{}, false
	}
	sorted := SortByRevision(records)
	return sorted[0], true
}
func Summarize(records []domain.Record) map[domain.Status]int {
	out := map[domain.Status]int{}
	for _, r := range records {
		out[r.Status]++
	}
	return out
}
func IsArchived(records []domain.Record) bool {
	for _, r := range records {
		if r.Status == domain.StatusArchived {
			return true
		}
	}
	return false
}
