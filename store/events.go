package store

import (
	"gestureflame/domain"
	"go.etcd.io/bbolt"
)

func (d *DB) SaveEvent(e domain.AuditEvent) error {
	data, err := domain.EncodeAuditEvent(e)
	if err != nil {
		return err
	}
	return d.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucketFor("event")).Put([]byte(e.ID), data) })
}
func (d *DB) ListEvents(recordID string) ([]domain.AuditEvent, error) {
	out := []domain.AuditEvent{}
	err := d.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketFor("event")).ForEach(func(_, value []byte) error {
			e, err := domain.DecodeAuditEvent(append([]byte(nil), value...))
			if err == nil && (recordID == "" || e.RecordID == recordID) {
				out = append(out, e)
			}
			return err
		})
	})
	return out, err
}
func (d *DB) SaveWorkflow(w domain.Workflow) error {
	data, err := domain.EncodeWorkflow(w)
	if err != nil {
		return err
	}
	return d.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucketFor("workflow")).Put([]byte(w.ID), data) })
}
func (d *DB) GetWorkflow(id string) (domain.Workflow, error) {
	var out domain.Workflow
	err := d.View(func(tx *bbolt.Tx) error {
		data := tx.Bucket(bucketFor("workflow")).Get([]byte(id))
		if data == nil {
			return domain.ErrNotFound
		}
		var err error
		out, err = domain.DecodeWorkflow(append([]byte(nil), data...))
		return err
	})
	return out, err
}
func (d *DB) SaveAttachment(a domain.Attachment) error {
	data, err := domain.EncodeAttachment(a)
	if err != nil {
		return err
	}
	return d.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucketFor("attachment")).Put([]byte(a.ID), data) })
}
func (d *DB) ListAttachments(recordID string) ([]domain.Attachment, error) {
	out := []domain.Attachment{}
	err := d.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketFor("attachment")).ForEach(func(_, value []byte) error {
			a, err := domain.DecodeAttachment(append([]byte(nil), value...))
			if err == nil && (recordID == "" || a.RecordID == recordID) {
				out = append(out, a)
			}
			return err
		})
	})
	return out, err
}
func (d *DB) CountEvents(recordID string) (int, error) {
	all, err := d.ListEvents(recordID)
	return len(all), err
}
