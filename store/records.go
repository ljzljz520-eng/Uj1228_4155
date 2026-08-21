package store

import (
	"bytes"
	"gestureflame/domain"
	"go.etcd.io/bbolt"
	"sort"
)

func (d *DB) SaveRecord(r domain.Record) error {
	data, err := domain.EncodeRecord(r)
	if err != nil {
		return err
	}
	return d.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucketFor("record")).Put([]byte(r.ID), data) })
}
func (d *DB) GetRecord(id string) (domain.Record, error) {
	var out domain.Record
	err := d.View(func(tx *bbolt.Tx) error {
		data := tx.Bucket(bucketFor("record")).Get([]byte(id))
		if data == nil {
			return domain.ErrNotFound
		}
		var err error
		out, err = domain.DecodeRecord(append([]byte(nil), data...))
		return err
	})
	return out, err
}
func (d *DB) ListRecords() ([]domain.Record, error) {
	out := []domain.Record{}
	err := d.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketFor("record")).ForEach(func(_, value []byte) error {
			r, err := domain.DecodeRecord(append([]byte(nil), value...))
			if err == nil {
				out = append(out, r)
			}
			return err
		})
	})
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, err
}
func (d *DB) FindRecords(query string) ([]domain.Record, error) {
	all, err := d.ListRecords()
	if err != nil {
		return nil, err
	}
	q := bytes.ToLower([]byte(query))
	out := make([]domain.Record, 0)
	for _, r := range all {
		if bytes.Contains(bytes.ToLower([]byte(r.SearchText())), q) {
			out = append(out, r)
		}
	}
	return out, nil
}
func (d *DB) DeleteRecord(id string) error {
	return d.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucketFor("record")).Delete([]byte(id)) })
}
func (d *DB) CountRecords() (int, error) { all, err := d.ListRecords(); return len(all), err }
func (d *DB) SaveMany(records []domain.Record) error {
	return d.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketFor("record"))
		for _, r := range records {
			data, err := domain.EncodeRecord(r)
			if err != nil {
				return err
			}
			if err := b.Put([]byte(r.ID), data); err != nil {
				return err
			}
		}
		return nil
	})
}
