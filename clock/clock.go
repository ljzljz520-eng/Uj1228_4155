package clock

import "fmt"

type Clock struct{ tick int64 }

func New() *Clock { return &Clock{tick: 1} }
func (c *Clock) Now() string {
	value := fmt.Sprintf("2026-01-01T00:00:%02dZ", c.tick%60)
	c.tick++
	return value
}
func (c *Clock) Peek() string               { return fmt.Sprintf("2026-01-01T00:00:%02dZ", c.tick%60) }
func (c *Clock) Reset()                     { c.tick = 1 }
func (c *Clock) Stamp(prefix string) string { return fmt.Sprintf("%s-%04d", prefix, c.tick) }
