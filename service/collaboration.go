package service

import (
	"errors"
	"gestureflame/domain"
)

func (s *Service) Search(query string) ([]domain.Record, error) { return s.db.FindRecords(query) }
func (s *Service) Collaborate(id, actor, result string, expectedRevision int) (domain.Record, error) {
	if err := s.ensureActor(actor); err != nil {
		return domain.Record{}, err
	}
	r, err := s.load(id)
	if err != nil {
		return r, err
	}
	if expectedRevision != r.Revision {
		return r, errors.New("revision conflict")
	}
	if r.IsTerminal() && r.Status != domain.StatusArchived {
		return r, errors.New("published record cannot change")
	}
	if result == "" {
		return r, domain.ErrMissingField
	}
	r.Result = domain.NormalizeResult(result)
	r.Revision++
	if err := s.update(r); err != nil {
		return r, err
	}
	if err := s.recordAudit(r, "collaborate", actor, "result revised"); err != nil {
		return r, err
	}
	return r, nil
}
func (s *Service) PublishWithConfirmation(id, actor string, expected string) (domain.Record, error) {
	r, err := s.load(id)
	if err != nil {
		return r, err
	}
	if !domain.SameResult(r.Result, expected) {
		return r, errors.New("confirmation mismatch")
	}
	return s.Publish(id, actor)
}
func (s *Service) History(id string) ([]domain.AuditEvent, error) { return s.db.ListEvents(id) }
func (s *Service) IsIndependentState(id string, prior string) (bool, error) {
	r, err := s.load(id)
	if err != nil {
		return false, err
	}
	return !domain.SameResult(r.Result, prior), nil
}
func (s *Service) Count(query string) (int, error) { rs, err := s.Search(query); return len(rs), err }
