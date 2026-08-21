package service

import (
	"errors"
	"gestureflame/domain"
)

type ReviewResult struct {
	Record    domain.Record
	Checklist domain.ReviewChecklist
	Issues    []domain.FieldIssue
}

func (s *Service) BuildChecklist(id string) (ReviewResult, error) {
	r, err := s.load(id)
	if err != nil {
		return ReviewResult{}, err
	}
	issues := domain.ValidateFields(r)
	checklist := domain.ReviewChecklist{Identity: r.BatchCode != "" && r.Title != "", Evidence: r.Revision > 0, Result: r.Result != "", Ownership: len(r.CreatedAt) > 0}
	return ReviewResult{Record: r, Checklist: checklist, Issues: issues}, nil
}
func (s *Service) ApproveChecklist(id, actor string) (domain.Record, error) {
	if !domain.ValidateActor(actor) {
		return domain.Record{}, errors.New("invalid actor")
	}
	result, err := s.BuildChecklist(id)
	if err != nil {
		return domain.Record{}, err
	}
	if !result.Checklist.Complete() {
		return result.Record, errors.New(domain.ExplainIssues(result.Issues))
	}
	return s.Review(id, actor)
}
func (s *Service) AddEvidence(id, actor, detail string) error {
	if !domain.ValidateActor(actor) || detail == "" {
		return errors.New("invalid evidence")
	}
	r, err := s.load(id)
	if err != nil {
		return err
	}
	return s.recordAudit(r, "evidence", actor, detail)
}
func (s *Service) ConfirmIndependent(id, actor, expected string) (domain.Record, error) {
	r, err := s.load(id)
	if err != nil {
		return r, err
	}
	if !domain.SameResult(r.Result, expected) {
		return r, errors.New("independent result mismatch")
	}
	return s.Confirm(id, actor)
}
func (s *Service) ArchiveWithEvidence(id, actor string) (domain.Record, error) {
	if err := s.AddEvidence(id, actor, "archive evidence attached"); err != nil {
		return domain.Record{}, err
	}
	return s.Archive(id, actor)
}
func (s *Service) ReopenCheck(id string) (bool, error) {
	_, err := s.db.GetRecord(id)
	return err == nil, err
}
func (s *Service) HasAuditAction(id, action string) (bool, error) {
	es, err := s.History(id)
	if err != nil {
		return false, err
	}
	for _, e := range es {
		if e.Action == action {
			return true, nil
		}
	}
	return false, nil
}
func (s *Service) AuditCount(id string) (int, error) { return s.db.CountEvents(id) }
