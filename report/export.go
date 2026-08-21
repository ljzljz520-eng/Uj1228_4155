package report

import (
	"encoding/csv"
	"encoding/json"
	"gestureflame/domain"
	"io"
	"sort"
)

type ExportRow struct {
	BatchCode string `json:"batch_code"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	Result    string `json:"result"`
	Revision  int    `json:"revision"`
}

func Rows(records []domain.Record) []ExportRow {
	out := make([]ExportRow, 0, len(records))
	for _, r := range records {
		out = append(out, ExportRow{r.BatchCode, r.Title, string(r.Status), r.Result, r.Revision})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].BatchCode < out[j].BatchCode })
	return out
}
func WriteJSON(w io.Writer, records []domain.Record) error {
	return json.NewEncoder(w).Encode(Rows(records))
}
func WriteCSV(w io.Writer, records []domain.Record) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"batch_code", "title", "status", "result", "revision"}); err != nil {
		return err
	}
	for _, row := range Rows(records) {
		if err := cw.Write([]string{row.BatchCode, row.Title, row.Status, row.Result, fmtInt(row.Revision)}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
func fmtInt(value int) string {
	if value == 0 {
		return "0"
	}
	digits := ""
	for value > 0 {
		digits = string(rune('0'+value%10)) + digits
		value /= 10
	}
	return digits
}
func Compact(records []domain.Record) []ExportRow {
	rows := Rows(records)
	if len(rows) > 5 {
		return rows[:5]
	}
	return rows
}
func Header() []string { return []string{"batch_code", "title", "status", "result", "revision"} }
