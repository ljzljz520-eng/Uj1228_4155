package domain

import (
	"fmt"
	"strings"
)

func EventFor(r Record, id, action, actor, detail, at string) AuditEvent {
	return AuditEvent{ID: id, RecordID: r.ID, Action: action, Actor: actor, Detail: detail, CreatedAt: at}
}
func IsMutation(action string) bool {
	switch action {
	case "create", "review", "confirm", "archive", "publish", "collaborate":
		return true
	default:
		return false
	}
}
func IsRead(action string) bool {
	return action == "search" || action == "recall" || action == "report"
}
func EventSummary(e AuditEvent) string { return fmt.Sprintf("%s:%s:%s", e.Action, e.Actor, e.Detail) }
func FilterEvents(events []AuditEvent, action string) []AuditEvent {
	out := []AuditEvent{}
	for _, e := range events {
		if action == "" || e.Action == action {
			out = append(out, e)
		}
	}
	return out
}
func EventActors(events []AuditEvent) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, e := range events {
		actor := strings.TrimSpace(e.Actor)
		if actor != "" && !seen[actor] {
			seen[actor] = true
			out = append(out, actor)
		}
	}
	return out
}
func EventCountByAction(events []AuditEvent) map[string]int {
	out := map[string]int{}
	for _, e := range events {
		out[e.Action]++
	}
	return out
}
func HasMutation(events []AuditEvent) bool {
	for _, e := range events {
		if IsMutation(e.Action) {
			return true
		}
	}
	return false
}
func LastMutation(events []AuditEvent) (AuditEvent, bool) {
	for i := len(events) - 1; i >= 0; i-- {
		if IsMutation(events[i].Action) {
			return events[i], true
		}
	}
	return AuditEvent{}, false
}
