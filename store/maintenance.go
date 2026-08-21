package store

import (
	"encoding/json"
	"errors"
	"gestureflame/domain"
	"go.etcd.io/bbolt"
	"io"
	"sort"
)

func (d *DB) Export(w io.Writer) error {
	snapshot, err := d.Snapshot()
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(snapshot)
}
func (d *DB) Import(data []byte) error {
	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return err
	}
	if len(snapshot.Records) == 0 {
		return errors.New("snapshot has no records")
	}
	if err := d.SaveMany(snapshot.Records); err != nil {
		return err
	}
	for _, e := range snapshot.Events {
		if err := d.SaveEvent(e); err != nil {
			return err
		}
	}
	for _, w := range snapshot.Workflows {
		if err := d.SaveWorkflow(w); err != nil {
			return err
		}
	}
	for _, a := range snapshot.Attachments {
		if err := d.SaveAttachment(a); err != nil {
			return err
		}
	}
	return nil
}
func (d *DB) ValidateSnapshot(s Snapshot) []string {
	issues := []string{}
	ids := map[string]bool{}
	for _, r := range s.Records {
		if ids[r.ID] {
			issues = append(issues, "duplicate record "+r.ID)
		}
		ids[r.ID] = true
		if err := domain.ValidateRecord(r); err != nil {
			issues = append(issues, r.ID+": "+err.Error())
		}
	}
	for _, e := range s.Events {
		if !ids[e.RecordID] {
			issues = append(issues, "orphan event "+e.ID)
		}
	}
	for _, a := range s.Attachments {
		if !ids[a.RecordID] {
			issues = append(issues, "orphan attachment "+a.ID)
		}
	}
	return issues
}
func (d *DB) IDs() ([]string, error) {
	rs, err := d.ListRecords()
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rs))
	for _, r := range rs {
		out = append(out, r.ID)
	}
	sort.Strings(out)
	return out, nil
}
func (d *DB) HasRecord(id string) bool { _, err := d.GetRecord(id); return err == nil }
func (d *DB) EntityCounts() (map[string]int, error) {
	rs, err := d.ListRecords()
	if err != nil {
		return nil, err
	}
	es, err := d.ListEvents("")
	if err != nil {
		return nil, err
	}
	as, err := d.ListAttachments("")
	if err != nil {
		return nil, err
	}
	snapshot, err := d.Snapshot()
	if err != nil {
		return nil, err
	}
	return map[string]int{"records": len(rs), "events": len(es), "workflows": len(snapshot.Workflows), "attachments": len(as)}, nil
}
func (d *DB) Empty() error {
	return d.Update(func(tx *bbolt.Tx) error {
		for _, name := range []string{"records", "events", "workflows", "attachments"} {
			b := tx.Bucket([]byte(name))
			keys := [][]byte{}
			if err := b.ForEach(func(k, v []byte) error { keys = append(keys, append([]byte(nil), k...)); return nil }); err != nil {
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
