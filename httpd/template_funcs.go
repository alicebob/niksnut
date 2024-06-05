package httpd

import (
	"fmt"
	"html/template"
)

func showerror(e string) template.HTML {
	if e == "" {
		return ""
	}
	return template.HTML(fmt.Sprintf(
		`
                        <h2 class="error">Error</h2>
                        <span class="error">- %s</span><br />
                        <br />
                        <br />
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
