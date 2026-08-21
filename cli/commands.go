package cli

import (
	"errors"
	"strings"
)

type Command struct {
	Name string
	Args []string
}

func Parse(args []string) (Command, error) {
	if len(args) == 0 {
		return Command{}, errors.New("command required")
	}
	name := strings.TrimSpace(args[0])
	if name == "" {
		return Command{}, errors.New("command required")
	}
	return Command{Name: name, Args: append([]string(nil), args[1:]...)}, nil
}
func (c Command) WantsJSON() bool {
	for _, arg := range c.Args {
		if arg == "--json" {
			return true
		}
	}
	return false
}
func (c Command) Argument(index int) string {
	if index < 0 || index >= len(c.Args) {
		return ""
	}
	return c.Args[index]
}
func (c Command) Valid() bool {
	switch c.Name {
	case "import", "search", "workflow":
		return true
	default:
		return false
	}
}
func JoinArguments(args []string) string { return strings.Join(args, " ") }
func StripFlags(args []string) []string {
	out := []string{}
	for _, arg := range args {
		if !strings.HasPrefix(arg, "--") {
			out = append(out, arg)
		}
	}
	return out
}
