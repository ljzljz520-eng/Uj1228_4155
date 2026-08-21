package service

import (
	"gestureflame/domain"
	"sort"
)

type Metrics struct {
	Total     int
	Draft     int
	Reviewed  int
	Confirmed int
	Archived  int
	Published int
	Revisions int
}

func CalculateMetrics(records []domain.Record) Metrics {
	m := Metrics{}
	m.Total = len(records)
	for _, r := range records {
		m.Revisions += r.Revision
		switch r.Status {
		case domain.StatusDraft:
			m.Draft++
		case domain.StatusReviewed:
			m.Reviewed++
		case domain.StatusConfirmed:
			m.Confirmed++
		case domain.StatusArchived:
			m.Archived++
		case domain.StatusPublished:
			m.Published++
		}
	}
	return m
}
func (s *Service) Metrics(query string) (Metrics, error) {
	rs, err := s.Search(query)
	return CalculateMetrics(rs), err
}
func (s *Service) PendingReview() ([]domain.Record, error) {
	return s.db.SearchByStatus(domain.StatusDraft)
}
func (s *Service) ConfirmedBatches() ([]domain.Record, error) {
	return s.db.SearchByStatus(domain.StatusConfirmed)
}
func (s *Service) ArchiveQueue() ([]domain.Record, error) {
	return s.db.SearchByStatus(domain.StatusArchived)
}
func SortRecords(records []domain.Record) []domain.Record {
	out := append([]domain.Record(nil), records...)
	sort.Slice(out, func(i, j int) bool {
		if out[i].BatchCode == out[j].BatchCode {
			return out[i].Revision < out[j].Revision
		}
		return out[i].BatchCode < out[j].BatchCode
	})
	return out
}
func GroupByBatch(records []domain.Record) map[string][]domain.Record {
	out := map[string][]domain.Record{}
	for _, r := range records {
		out[r.BatchCode] = append(out[r.BatchCode], r)
	}
	return out
}
