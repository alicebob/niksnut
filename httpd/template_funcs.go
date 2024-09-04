package httpd

import (
	"fmt"
	"html/template"
	"time"
)

var (
	BrowserTimeFormat = "2006-01-02T15:04:05"
)

func showerror(e string) template.HTML {
	if e == "" {
		return ""
	}
	return template.HTML(fmt.Sprintf(
		`
			<h2>Error</h2>
			<div class="col col1">
                        Error: %s
			</div>
                    `,
		template.HTMLEscapeString(e),
	))
}

func htmlRadio(name, value, current string, label string) template.HTML {
	c := ""
	if value == current {
		c = " CHECKED"
	}
	return template.HTML(fmt.Sprintf(
		`<label><input type="radio" name="%s" value="%s"%s /> %s</label>`,
		template.HTMLEscapeString(name),
		template.HTMLEscapeString(value),
		c,
		template.HTMLEscapeString(label),
	))
}

// little colored dot with a label
// status can be "running", "success", "failed".
func showstatus(status, label string) template.HTML {
	switch status {
	case "running",
		"success",
		"failed":
	default:
		panic("no my compiler can't check this")
	}

	return template.HTML(fmt.Sprintf(
		`<div class="status %[1]s" title="%[2]s"><div class="icon"></div> %[2]s</div>`,
		status,
		template.HTMLEscapeString(label),
	))
}

func datetime(t time.Time) template.HTML {
	if t.IsZero() {
		return "-"
	}
	return template.HTML(fmt.Sprintf(`<span onmouseover="hoverLocaltime(this)">%s</span>`, t.UTC().Format(BrowserTimeFormat)))
}

func duration(t time.Duration) string {
	if t <= 0 {
		return ""
	}
	// I prefer to show the same unit everywhere, easier to compare visually.
	return fmt.Sprintf("%.02fs", t.Seconds())
	// return t.Truncate(10 * time.Millisecond).String()
}
