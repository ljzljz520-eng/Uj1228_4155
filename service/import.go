package service

import (
	"encoding/json"
	"gestureflame/domain"
)

func ParseImport(data []byte) (domain.ImportDraft, error) {
	var d domain.ImportDraft
	if err := json.Unmarshal(data, &d); err != nil {
		return d, err
	}
	return d, domain.ValidateImport(d)
}
func (s *Service) ImportJSON(data []byte) (domain.Record, error) {
	d, err := ParseImport(data)
	if err != nil {
		return domain.Record{}, err
	}
	return s.Import(d)
}
func (s *Service) ReportData(id string) (domain.Record, []domain.AuditEvent, []domain.Attachment, error) {
	r, e, err := s.Recall(id)
	if err != nil {
		return r, nil, nil, err
	}
	a, err := s.db.ListAttachments(id)
	return r, e, a, err
}
func (s *Service) RegisterAttachment(id string, d domain.ImportDraft) error {
	r, err := s.load(id)
	if err != nil {
		return err
	}
	return s.recordAttachment(r, d)
}
func (s *Service) StartWorkflow(id, owner string) (domain.Workflow, error) {
	r, err := s.load(id)
	if err != nil {
		return domain.Workflow{}, err
	}
	w := domain.Workflow{ID: s.nextID("workflow"), RecordID: r.ID, Name: "handoff", Stage: "started", Owner: owner, StartedAt: s.clock.Now()}
	return w, s.db.SaveWorkflow(w)
}
func (s *Service) CompleteWorkflow(w domain.Workflow) (domain.Workflow, error) {
	if w.Stage == "completed" {
		return w, nil
	}
	w.Stage = "completed"
	w.CompletedAt = s.clock.Now()
	return w, s.db.SaveWorkflow(w)
}
