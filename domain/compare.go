package domain

import "strings"

type Difference struct {
	Field  string
	Before string
	After  string
}

func CompareRecords(before, after Record) []Difference {
	out := []Difference{}
	if before.BatchCode != after.BatchCode {
		out = append(out, Difference{"batch_code", before.BatchCode, after.BatchCode})
	}
	if before.Title != after.Title {
		out = append(out, Difference{"title", before.Title, after.Title})
	}
	if before.Status != after.Status {
		out = append(out, Difference{"status", string(before.Status), string(after.Status)})
	}
	if before.Result != after.Result {
		out = append(out, Difference{"result", before.Result, after.Result})
	}
	if before.Revision != after.Revision {
		out = append(out, Difference{"revision", itoa(before.Revision), itoa(after.Revision)})
	}
	return out
}
func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	sign := ""
	if value < 0 {
		sign = "-"
		value = -value
	}
	out := ""
	for value > 0 {
		out = string(rune('0'+value%10)) + out
		value /= 10
	}
	return sign + out
}
func DifferenceNames(diffs []Difference) []string {
	out := []string{}
	for _, d := range diffs {
		out = append(out, d.Field)
	}
	return out
}
func HasDifference(diffs []Difference, field string) bool {
	for _, d := range diffs {
		if d.Field == field {
			return true
		}
	}
	return false
}
func NormalizeDifferences(diffs []Difference) []Difference {
	out := make([]Difference, len(diffs))
	copy(out, diffs)
	for i := range out {
		out[i].Before = strings.TrimSpace(out[i].Before)
		out[i].After = strings.TrimSpace(out[i].After)
	}
	return out
}
func EqualIgnoringRevision(a, b Record) bool {
	a.Revision = b.Revision
	return len(CompareRecords(a, b)) == 0
}
