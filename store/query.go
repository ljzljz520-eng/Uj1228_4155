package store

import (
	"encoding/json"
	"gestureflame/domain"
	"go.etcd.io/bbolt"
	"sort"
	"strings"
)

type Snapshot struct {
	Records     []domain.Record     `json:"records"`
	Events      []domain.AuditEvent `json:"events"`
	Workflows   []domain.Workflow   `json:"workflows"`
	Attachments []domain.Attachment `json:"attachments"`
}

func (d *DB) Snapshot() (Snapshot, error) {
	var s Snapshot
	rs, err := d.ListRecords()
	if err != nil {
		return s, err
	}
	s.Records = rs
	es, err := d.ListEvents("")
	if err != nil {
		return s, err
	}
	s.Events = es
	as, err := d.ListAttachments("")
	if err != nil {
		return s, err
	}
	s.Attachments = as
	err = d.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketFor("workflow")).ForEach(func(_, value []byte) error {
			w, e := domain.DecodeWorkflow(append([]byte(nil), value...))
			if e == nil {
				s.Workflows = append(s.Workflows, w)
			}
			return e
		})
	})
	sort.Slice(s.Workflows, func(i, j int) bool { return s.Workflows[i].ID < s.Workflows[j].ID })
	return s, err
}
func (d *DB) SnapshotJSON() ([]byte, error) {
	s, err := d.Snapshot()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(s, "", "  ")
}
func (d *DB) SearchByStatus(status domain.Status) ([]domain.Record, error) {
	rs, err := d.ListRecords()
	if err != nil {
		return nil, err
	}
	out := []domain.Record{}
	for _, r := range rs {
		if r.Status == status {
			out = append(out, r)
		}
	}
	return out, nil
}
func (d *DB) SearchByPrefix(prefix string) ([]domain.Record, error) {
	rs, err := d.ListRecords()
	if err != nil {
		return nil, err
	}
	p := strings.ToUpper(prefix)
	out := []domain.Record{}
	for _, r := range rs {
		if strings.HasPrefix(r.BatchCode, p) {
			out = append(out, r)
		}
	}
	return out, nil
}
func (d *DB) LatestEvent(recordID string) (domain.AuditEvent, error) {
	es, err := d.ListEvents(recordID)
	if err != nil {
		return domain.AuditEvent{}, err
	}
	if len(es) == 0 {
		return domain.AuditEvent{}, domain.ErrNotFound
	}
	sort.Slice(es, func(i, j int) bool { return es[i].ID > es[j].ID })
	return es[0], nil
}
func (d *DB) WorkflowsFor(recordID string) ([]domain.Workflow, error) {
	var out []domain.Workflow
	err := d.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketFor("workflow")).ForEach(func(_, value []byte) error {
			w, e := domain.DecodeWorkflow(append([]byte(nil), value...))
			if e == nil && w.RecordID == recordID {
				out = append(out, w)
			}
			return e
		})
	})
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, err
}
func (d *DB) PurgeRecord(id string) error {
	return d.Update(func(tx *bbolt.Tx) error {
		for _, kind := range []string{"record", "event", "workflow", "attachment"} {
			b := tx.Bucket(bucketFor(kind))
			if kind == "record" {
				if err := b.Delete([]byte(id)); err != nil {
					return err
				}
			}
			var keys [][]byte
			if err := b.ForEach(func(k, v []byte) error {
				if kind != "record" && bytesEqualField(v, []byte(id)) {
					keys = append(keys, append([]byte(nil), k...))
				}
				return nil
			}); err != nil {
				return err
			}
			for _, k := range keys {
				if err := b.Delete(k); err != nil {
					return err
				}
			}
		}
		return nil
	})
}
func bytesEqualField(data, id []byte) bool {
	return len(data) > 0 && strings.Contains(string(data), string(id))
}
