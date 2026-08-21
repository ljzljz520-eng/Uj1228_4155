package service

import (
	"errors"
	"fmt"
	"gestureflame/clock"
	"gestureflame/domain"
	"gestureflame/store"
)

type ArchiveCallback func(domain.Record) error

type Service struct {
	db              *store.DB
	clock           *clock.Clock
	archiveCallback ArchiveCallback
}

func New(db *store.DB, c *clock.Clock) *Service          { return &Service{db: db, clock: c} }
func (s *Service) SetArchiveCallback(cb ArchiveCallback) { s.archiveCallback = cb }
func (s *Service) DB() *store.DB                         { return s.db }
// runArchiveCallback invokes the archive callback, treating a missing callback as
// a no-op so that archiving a batch never depends on external callback state being
// configured. Each batch must archive independently and preserve its own confirmed
// result, even when a second identical processing is triggered without a callback.
func (s *Service) runArchiveCallback(r domain.Record) error {
	if s.archiveCallback == nil {
		return nil
	}
	return s.archiveCallback(r)
}
func (s *Service) nextID(prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, s.clock.Stamp(prefix))
}
func (s *Service) makeRecord(d domain.ImportDraft) domain.Record {
	now := s.clock.Now()
	return domain.Record{ID: s.nextID("record"), BatchCode: domain.NormalizeBatchCode(d.BatchCode), Title: domain.NormalizeTitle(d.Title), Status: domain.StatusDraft, Result: domain.NormalizeResult(d.Result), Revision: 1, CreatedAt: now, UpdatedAt: now}
}
func (s *Service) Create(d domain.ImportDraft) (domain.Record, error) {
	if err := domain.ValidateImport(d); err != nil {
		return domain.Record{}, err
	}
	r := s.makeRecord(d)
	if err := s.db.SaveRecord(r); err != nil {
		return domain.Record{}, err
	}
	if err := s.recordAudit(r, "create", d.Actor, "record created"); err != nil {
		return domain.Record{}, err
	}
	return r, nil
}
func (s *Service) recordAudit(r domain.Record, action, actor, detail string) error {
	e := domain.AuditEvent{ID: s.nextID("event"), RecordID: r.ID, Action: action, Actor: actor, Detail: detail, CreatedAt: s.clock.Now()}
	return s.db.SaveEvent(e)
}
func (s *Service) recordWorkflow(r domain.Record, name, stage, owner string) error {
	w := domain.Workflow{ID: s.nextID("workflow"), RecordID: r.ID, Name: name, Stage: stage, Owner: owner, StartedAt: r.CreatedAt, CompletedAt: s.clock.Now()}
	return s.db.SaveWorkflow(w)
}
func (s *Service) recordAttachment(r domain.Record, d domain.ImportDraft) error {
	a := domain.Attachment{ID: s.nextID("attachment"), RecordID: r.ID, Name: d.AttachmentName, MediaType: d.AttachmentType, Digest: d.AttachmentDigest, CreatedAt: s.clock.Now()}
	return s.db.SaveAttachment(a)
}
func (s *Service) load(id string) (domain.Record, error) {
	r, err := s.db.GetRecord(id)
	if err != nil {
		return r, err
	}
	if err := domain.ValidateRecord(r); err != nil {
		return r, err
	}
	return r, nil
}
func (s *Service) update(r domain.Record) error {
	r.UpdatedAt = s.clock.Now()
	return s.db.SaveRecord(r)
}
func (s *Service) ensureActor(actor string) error {
	if actor == "" {
		return errors.New("actor required")
	}
	return nil
}
