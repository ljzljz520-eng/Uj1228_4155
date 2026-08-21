package service

import (
	"errors"
	"gestureflame/domain"
	"sort"
	"strings"
)

type BatchOutcome struct {
	Record      domain.Record
	Events      int
	Attachments int
	Ready       bool
}

func (s *Service) Inspect(id string) (BatchOutcome, error) {
	r, e, a, err := s.ReportData(id)
	if err != nil {
		return BatchOutcome{}, err
	}
	ready := r.Status == domain.StatusConfirmed || r.Status == domain.StatusArchived || r.Status == domain.StatusPublished
	return BatchOutcome{r, len(e), len(a), ready}, nil
}
func (s *Service) EnsureArchiveReady(id string) error {
	r, err := s.load(id)
	if err != nil {
		return err
	}
	if !domain.CanArchive(r) {
		return errors.New("record is not confirmed")
	}
	ok, err := s.HasAuditAction(id, "confirm")
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("confirmation evidence missing")
	}
	return nil
}
func (s *Service) ArchiveVerified(id, actor string, cb ArchiveCallback) (domain.Record, error) {
	if err := s.EnsureArchiveReady(id); err != nil {
		return domain.Record{}, err
	}
	previous := s.archiveCallback
	s.archiveCallback = cb
	r, err := s.Archive(id, actor)
	s.archiveCallback = previous
	return r, err
}
func (s *Service) Reconcile(id, actor, expected string) (domain.Record, error) {
	r, err := s.load(id)
	if err != nil {
		return r, err
	}
	if !domain.SameResult(r.Result, expected) {
		if _, e := s.Collaborate(id, actor, expected, r.Revision); e != nil {
			return r, e
		}
		r, err = s.load(id)
	}
	return r, err
}
func (s *Service) ReconcileAndPublish(id, actor, expected string) (domain.Record, error) {
	r, err := s.Reconcile(id, actor, expected)
	if err != nil {
		return r, err
	}
	if r.Status == domain.StatusArchived {
		return s.PublishWithConfirmation(id, actor, expected)
	}
	if r.Status == domain.StatusPublished {
		return r, nil
	}
	return r, errors.New("record not ready for publish")
}
func (s *Service) BatchImport(drafts []domain.ImportDraft) ([]domain.Record, error) {
	out := make([]domain.Record, 0, len(drafts))
	for _, d := range drafts {
		r, err := s.Import(d)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}
func (s *Service) ValidateBatch(records []domain.Record) []domain.FieldIssue {
	issues := []domain.FieldIssue{}
	seen := map[string]bool{}
	for _, r := range records {
		if seen[r.BatchCode] {
			issues = append(issues, domain.FieldIssue{"batch_code", "duplicate " + r.BatchCode})
		}
		seen[r.BatchCode] = true
		issues = append(issues, domain.ValidateFields(r)...)
	}
	return issues
}
func (s *Service) ResolveDuplicates(records []domain.Record) []domain.Record {
	grouped := map[string]domain.Record{}
	for _, r := range records {
		current, ok := grouped[r.BatchCode]
		if !ok || r.Revision > current.Revision {
			grouped[r.BatchCode] = r
		}
	}
	out := make([]domain.Record, 0, len(grouped))
	for _, r := range grouped {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].BatchCode < out[j].BatchCode })
	return out
}
func (s *Service) FindByActor(actor string) ([]domain.Record, error) {
	all, err := s.Search("")
	if err != nil {
		return nil, err
	}
	events, err := s.db.ListEvents("")
	if err != nil {
		return nil, err
	}
	ids := map[string]bool{}
	for _, e := range events {
		if strings.EqualFold(e.Actor, actor) {
			ids[e.RecordID] = true
		}
	}
	out := []domain.Record{}
	for _, r := range all {
		if ids[r.ID] {
			out = append(out, r)
		}
	}
	return out, nil
}
func (s *Service) ArchiveCount() (int, error) { rs, err := s.ArchiveQueue(); return len(rs), err }
func (s *Service) PublishedCount() (int, error) {
	rs, err := s.db.SearchByStatus(domain.StatusPublished)
	return len(rs), err
}
