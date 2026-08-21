package cli

import (
	"encoding/json"
	"fmt"
	"gestureflame/domain"
	"gestureflame/report"
	"gestureflame/service"
	"io"
	"strings"
)

func Run(app *service.Service, args []string, out io.Writer) error {
	if len(args) == 0 {
		_, err := fmt.Fprintln(out, Usage())
		return err
	}
	switch args[0] {
	case "import":
		return runImport(app, args[1:], out)
	case "search":
		return runSearch(app, args[1:], out)
	case "workflow":
		return runWorkflow(app, args[1:], out)
	default:
		return fmt.Errorf("unknown command %s", args[0])
	}
}
func runImport(app *service.Service, args []string, out io.Writer) error {
	if len(args) != 1 {
		return fmt.Errorf("import expects JSON")
	}
	r, err := app.ImportJSON([]byte(args[0]))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, r.ID)
	return err
}
func runSearch(app *service.Service, args []string, out io.Writer) error {
	q := strings.Join(args, " ")
	rs, err := app.Search(q)
	if err != nil {
		return err
	}
	if len(rs) == 0 {
		_, err = io.WriteString(out, report.EmptyReport())
		return err
	}
	_, err = io.WriteString(out, report.RenderList(rs))
	return err
}
func runWorkflow(app *service.Service, args []string, out io.Writer) error {
	if len(args) < 2 {
		return fmt.Errorf("workflow expects id owner")
	}
	w, err := app.StartWorkflow(args[0], args[1])
	if err != nil {
		return err
	}
	data, err := json.Marshal(w)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, string(data))
	return err
}
func ParseDraft(data string) (domain.ImportDraft, error) {
	var d domain.ImportDraft
	err := json.Unmarshal([]byte(data), &d)
	return d, err
}
func Usage() string {
	return "flamectl import '<json>' | search <query> | workflow <record-id> <owner>"
}
