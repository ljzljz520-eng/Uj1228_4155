package domain

import (
	"errors"
	"strings"
)

type Status string

const (
	StatusDraft     Status = "draft"
	StatusReviewed  Status = "reviewed"
	StatusConfirmed Status = "confirmed"
	StatusArchived  Status = "archived"
	StatusPublished Status = "published"
)

type Record struct {
	ID        string `json:"id"`
	BatchCode string `json:"batch_code"`
	Title     string `json:"title"`
	Status    Status `json:"status"`
	Result    string `json:"result"`
	Revision  int    `json:"revision"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type AuditEvent struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	Action    string `json:"action"`
	Actor     string `json:"actor"`
	Detail    string `json:"detail"`
	CreatedAt string `json:"created_at"`
}

type Workflow struct {
	ID          string `json:"id"`
	RecordID    string `json:"record_id"`
	Name        string `json:"name"`
	Stage       string `json:"stage"`
	Owner       string `json:"owner"`
	StartedAt   string `json:"started_at"`
	CompletedAt string `json:"completed_at"`
}

type Attachment struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	Name      string `json:"name"`
	MediaType string `json:"media_type"`
	Digest    string `json:"digest"`
	CreatedAt string `json:"created_at"`
}

type ImportDraft struct {
	BatchCode        string `json:"batch_code"`
	Title            string `json:"title"`
	Result           string `json:"result"`
	Actor            string `json:"actor"`
	AttachmentName   string `json:"attachment_name"`
	AttachmentType   string `json:"attachment_type"`
	AttachmentDigest string `json:"attachment_digest"`
}

var ErrInvalidTransition = errors.New("invalid lifecycle transition")
var ErrMissingField = errors.New("required field missing")
var ErrNotFound = errors.New("record not found")

func (r Record) IsTerminal() bool { return r.Status == StatusArchived || r.Status == StatusPublished }
func (r Record) Clone() Record    { return r }
func (r Record) SearchText() string {
	return strings.ToLower(r.BatchCode + " " + r.Title + " " + r.Result)
}
