package renderer

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin/render"
)

// TemplateRenderer implements gin's HTMLRender interface
type TemplateRenderer struct {
	TemplatesPath string
}

// Instance renders a template with the given data
func (r *TemplateRenderer) Instance(name string, data interface{}) render.Render {
	return &templateRender{
		name:          name,
		data:          data,
		templatesPath: r.TemplatesPath,
	}
}

type templateRender struct {
	name          string
	data          interface{}
	templatesPath string
}

// Render writes the template to the response
func (t *templateRender) Render(w http.ResponseWriter) error {
	// Load base layout and the specific template
	baseTemplate := filepath.Join(t.templatesPath, "layouts", "base.html")
	pageTemplate := filepath.Join(t.templatesPath, t.name)

	tmpl, err := template.ParseFiles(baseTemplate, pageTemplate)
	if err != nil {
		// Log the error for debugging
		return fmt.Errorf("failed to parse templates (base: %s, page: %s): %w", baseTemplate, pageTemplate, err)
	}

	// Execute the "base" template
	return tmpl.ExecuteTemplate(w, "base", t.data)
}

// WriteContentType sets the Content-Type header
func (t *templateRender) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}
