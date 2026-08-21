package clock

import "fmt"

type Sequence struct{ value int }

func NewSequence() *Sequence { return &Sequence{value: 1} }
func (s *Sequence) Next(prefix string) string {
	value := fmt.Sprintf("%s-%04d", prefix, s.value)
	s.value++
	return value
}
func (s *Sequence) Current() int { return s.value }
func (s *Sequence) Reset()       { s.value = 1 }
func (s *Sequence) Batch(prefix string, count int) []string {
	out := []string{}
	for i := 0; i < count; i++ {
		out = append(out, s.Next(prefix))
	}
	return out
}
