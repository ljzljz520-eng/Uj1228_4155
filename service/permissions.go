package service

import (
	"errors"
	"gestureflame/domain"
	"strings"
)

type Role string

const (
	RoleOperator     Role = "operator"
	RoleReviewer     Role = "reviewer"
	RoleAuditor      Role = "auditor"
	RoleCollaborator Role = "collaborator"
)

func ParseRole(value string) Role {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "operator":
		return RoleOperator
	case "reviewer":
		return RoleReviewer
	case "auditor":
		return RoleAuditor
	case "collaborator":
		return RoleCollaborator
	default:
		return ""
	}
}
func CanPerform(role Role, action string) bool {
	switch action {
	case "create":
		return role == RoleOperator || role == RoleReviewer
	case "review", "confirm":
		return role == RoleReviewer
	case "archive", "publish":
		return role == RoleOperator || role == RoleReviewer
	case "collaborate":
		return role == RoleCollaborator || role == RoleReviewer
	case "report", "search":
		return role != ""
	default:
		return false
	}
}
func (s *Service) Authorize(actor, action string) error {
	parts := strings.SplitN(actor, ":", 2)
	role := ParseRole(parts[0])
	if !CanPerform(role, action) {
		return errors.New("permission denied")
	}
	return nil
}
func (s *Service) AuthorizedReview(id, actor string) (domain.Record, error) {
	if err := s.Authorize(actor, "review"); err != nil {
		return domain.Record{}, err
	}
	return s.Review(id, actor)
}
func (s *Service) AuthorizedArchive(id, actor string) (domain.Record, error) {
	if err := s.Authorize(actor, "archive"); err != nil {
		return domain.Record{}, err
	}
	return s.Archive(id, actor)
}
func (s *Service) AuthorizedCollaborate(id, actor, result string, revision int) (domain.Record, error) {
	if err := s.Authorize(actor, "collaborate"); err != nil {
		return domain.Record{}, err
	}
	return s.Collaborate(id, actor, result, revision)
}
func RequireOwner(owner, actor string) error {
	if owner == "" || actor == "" || owner != actor {
		return errors.New("owner mismatch")
	}
	return nil
}
func NormalizeActor(actor string) string { return strings.ToLower(strings.TrimSpace(actor)) }
func ActorRole(actor string) Role        { parts := strings.SplitN(actor, ":", 2); return ParseRole(parts[0]) }
