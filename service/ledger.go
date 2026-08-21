package service

import (
	"gestureflame/domain"
	"sort"
)

type LedgerEntry struct {
	Sequence int
	Event    domain.AuditEvent
}

func (s *Service) Ledger(id string) ([]LedgerEntry, error) {
	events, err := s.History(id)
	if err != nil {
		return nil, err
	}
	sort.Slice(events, func(i, j int) bool { return events[i].CreatedAt < events[j].CreatedAt })
	out := make([]LedgerEntry, len(events))
	for i, e := range events {
		out[i] = LedgerEntry{i + 1, e}
	}
	return out, nil
}
func (s *Service) LedgerActions(id string) ([]string, error) {
	ledger, err := s.Ledger(id)
	if err != nil {
		return nil, err
	}
	out := []string{}
	for _, entry := range ledger {
		out = append(out, entry.Event.Action)
	}
	return out, nil
}
func (s *Service) VerifyLedger(id string) error {
	r, events, err := s.Recall(id)
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return domain.ErrNotFound
	}
	if events[0].RecordID != r.ID {
		return domain.ErrNotFound
	}
	return nil
}
func (s *Service) HasIndependentConfirmation(id string) (bool, error) {
	events, err := s.History(id)
	if err != nil {
		return false, err
	}
	count := 0
	for _, e := range events {
		if e.Action == "confirm" || e.Action == "collaborate" {
			count++
		}
	}
	return count >= 1, nil
}
func (s *Service) EventDigest(id string) (string, error) {
	events, err := s.History(id)
	if err != nil {
		return "", err
	}
	digest := ""
	for _, e := range events {
		digest += e.Action + ":" + e.Detail + ";"
	}
	return digest, nil
}
func (s *Service) RequiresFollowup(id string) (bool, error) {
	r, err := s.load(id)
	if err != nil {
		return false, err
	}
	return r.Status == domain.StatusDraft || r.Status == domain.StatusReviewed, nil
}
func (s *Service) NextAction(id string) (string, error) {
	r, err := s.load(id)
	if err != nil {
		return "", err
	}
	return domain.ActionForStatus(r.Status), nil
}
