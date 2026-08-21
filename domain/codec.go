package domain

import "encoding/json"

func EncodeRecord(r Record) ([]byte, error) { return json.Marshal(r) }
func DecodeRecord(data []byte) (Record, error) {
	var r Record
	err := json.Unmarshal(data, &r)
	return r, err
}
func EncodeAuditEvent(e AuditEvent) ([]byte, error) { return json.Marshal(e) }
func DecodeAuditEvent(data []byte) (AuditEvent, error) {
	var e AuditEvent
	err := json.Unmarshal(data, &e)
	return e, err
}
func EncodeWorkflow(w Workflow) ([]byte, error) { return json.Marshal(w) }
func DecodeWorkflow(data []byte) (Workflow, error) {
	var w Workflow
	err := json.Unmarshal(data, &w)
	return w, err
}
func EncodeAttachment(a Attachment) ([]byte, error) { return json.Marshal(a) }
func DecodeAttachment(data []byte) (Attachment, error) {
	var a Attachment
	err := json.Unmarshal(data, &a)
	return a, err
}
