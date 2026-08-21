package clock

import "fmt"

func DayLabel(day int) string {
	if day < 1 {
		day = 1
	}
	if day > 31 {
		day = 31
	}
	return fmt.Sprintf("2026-01-%02d", day)
}
func SlotLabel(slot int) string {
	if slot < 0 {
		slot = 0
	}
	hour := slot % 24
	minute := (slot * 5) % 60
	return fmt.Sprintf("%02d:%02d", hour, minute)
}
func SequenceLabel(prefix string, value int) string {
	if value < 0 {
		value = 0
	}
	return fmt.Sprintf("%s-%04d", prefix, value)
}
func Window(start, end int) []string {
	if end < start {
		start, end = end, start
	}
	out := []string{}
	for i := start; i <= end; i++ {
		out = append(out, SlotLabel(i))
	}
	return out
}
