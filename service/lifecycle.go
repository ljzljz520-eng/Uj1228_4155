package service

import "gestureflame/domain"

func (s *Service) Review(id, actor string) (domain.Record, error) {
	if err := s.ensureActor(actor); err != nil {
		return domain.Record{}, err
	}
	r, err := s.load(id)
	if err != nil {
		return r, err
	}
	if !domain.CanReview(r) {
		return r, domain.ErrInvalidTransition
	}
	if err := domain.Transition(&r, domain.StatusReviewed); err != nil {
		return r, err
	}
	if err := s.update(r); err != nil {
		return r, err
	}
	if err := s.recordAudit(r, "review", actor, "fields reviewed"); err != nil {
		return r, err
	}
	return r, nil
}
func (s *Service) Confirm(id, actor string) (domain.Record, error) {
	if err := s.ensureActor(actor); err != nil {
		return domain.Record{}, err
	}
	r, err := s.load(id)
	if err != nil {
		return r, err
	}
	if !domain.CanConfirm(r) {
		return r, domain.ErrInvalidTransition
	}
	if err := domain.Transition(&r, domain.StatusConfirmed); err != nil {
		return r, err
	}
	if err := s.update(r); err != nil {
		return r, err
	}
	if err := s.recordAudit(r, "confirm", actor, "business result confirmed"); err != nil {
		return r, err
	}
	return r, nil
}
func (s *Service) Archive(id, actor string) (domain.Record, error) {
	if err := s.ensureActor(actor); err != nil {
		return domain.Record{}, err
	}
	r, err := s.load(id)
	if err != nil {
		return r, err
	}
	if !domain.CanArchive(r) {
		return r, domain.ErrInvalidTransition
	}
	if err := domain.Transition(&r, domain.StatusArchived); err != nil {
		return r, err
	}
	if err := s.update(r); err != nil {
		return r, err
	}
	if err := s.recordAudit(r, "archive", actor, "batch archived"); err != nil {
		return r, err
	}
	if err := s.runArchiveCallback(r); err != nil {
		return r, err
	}
	return r, nil
}
func (s *Service) Publish(id, actor string) (domain.Record, error) {
	if err := s.ensureActor(actor); err != nil {
		return domain.Record{}, err
	}
	r, err := s.load(id)
	if err != nil {
		return r, err
	}
	if !domain.CanPublish(r) {
		return r, domain.ErrInvalidTransition
	}
	if err := domain.Transition(&r, domain.StatusPublished); err != nil {
		return r, err
	}
	if err := s.update(r); err != nil {
		return r, err
	}
	if err := s.recordAudit(r, "publish", actor, "state published"); err != nil {
		return r, err
	}
	return r, nil
}
func (s *Service) Recall(id string) (domain.Record, []domain.AuditEvent, error) {
	r, err := s.load(id)
	if err != nil {
		return r, nil, err
	}
	e, err := s.db.ListEvents(id)
	return r, e, err
}
func (s *Service) RunCreateReviewConfirmArchive(d domain.ImportDraft, cb ArchiveCallback) (domain.Record, error) {
	r, err := s.Create(d)
	if err != nil {
		return r, err
	}
	if _, err = s.Review(r.ID, d.Actor); err != nil {
		return r, err
	}
	if _, err = s.Confirm(r.ID, d.Actor); err != nil {
		return r, err
	}
	s.SetArchiveCallback(cb)
	r, err = s.Archive(r.ID, d.Actor)
	return r, err
}
func (s *Service) Import(d domain.ImportDraft) (domain.Record, error) {
	r, err := s.Create(d)
	if err != nil {
		return r, err
	}
	if err := s.recordWorkflow(r, "import", "persisted", d.Actor); err != nil {
		return r, err
	}
	if err := s.recordAttachment(r, d); err != nil {
		return r, err
	}
	return r, nil
}
