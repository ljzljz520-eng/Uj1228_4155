package domain

func Transition(r *Record, next Status) error {
	if r == nil {
		return ErrMissingField
	}
	if !allowedTransition(r.Status, next) {
		return ErrInvalidTransition
	}
	r.Status = next
	r.Revision++
	return nil
}

func allowedTransition(from, to Status) bool {
	switch from {
	case StatusDraft:
		return to == StatusReviewed
	case StatusReviewed:
		return to == StatusConfirmed
	case StatusConfirmed:
		return to == StatusArchived
	case StatusArchived:
		return to == StatusPublished
	default:
		return false
	}
}

func CanReview(r Record) bool       { return r.Status == StatusDraft }
func CanConfirm(r Record) bool      { return r.Status == StatusReviewed && stringsNonEmpty(r.Result) }
func CanArchive(r Record) bool      { return r.Status == StatusConfirmed }
func CanPublish(r Record) bool      { return r.Status == StatusArchived }
func stringsNonEmpty(s string) bool { return len(s) > 0 }

func TransitionPath() []Status {
	return []Status{StatusDraft, StatusReviewed, StatusConfirmed, StatusArchived, StatusPublished}
}
func StatusRank(s Status) int {
	for i, v := range TransitionPath() {
		if s == v {
			return i
		}
	}
	return -1
}
