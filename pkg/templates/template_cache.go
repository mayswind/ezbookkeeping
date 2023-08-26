package templates

import (
	"fmt"
	"html/template"
	"path/filepath"
)

const templateBasePath = "templates"
const templateFileExtension = "tmpl"

var templateCache = make(map[string]*CachedTemplate)

// CachedTemplate represents a cached template
type CachedTemplate struct {
	templateName    string
	templateContent *template.Template
}

// GetTemplate returns a cached template instance according to the template name
func GetTemplate(templateName string) (*template.Template, error) {
	fullPath := filepath.Join(templateBasePath, fmt.Sprintf("%s.%s", templateName, templateFileExtension))

	cachedTemplate, exists := templateCache[templateName]

	if exists {
		return cachedTemplate.templateContent, nil
	}

	tmpl, err := template.ParseFiles(fullPath)

	if err != nil {
		return nil, err
	}

	templateCache[templateName] = &CachedTemplate{
		templateName:    templateName,
		templateContent: tmpl,
	}

	return tmpl, err
}
