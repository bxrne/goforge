package templates

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed templates
var templateFS embed.FS

type TemplateData struct {
	ProjectName string
	Template    string
	Database    string
	Auth        string
	GoVersion   string
}

type Engine struct {
	templates map[string]*template.Template
}

func NewEngine() (*Engine, error) {
	e := &Engine{
		templates: make(map[string]*template.Template),
	}

	// Load all template files
	templateFiles := []string{
		"templates/api/go.mod.tmpl",
		"templates/api/main.go.tmpl",
		"templates/api/env.tmpl",
		"templates/api/Dockerfile.tmpl",
		"templates/cli/go.mod.tmpl",
		"templates/cli/main.go.tmpl",
		"templates/cli/root.go.tmpl",
		"templates/cli/commands.go.tmpl",
		"templates/cli/README.md.tmpl",
		"templates/cli/Makefile.tmpl",
	}

	for _, templateFile := range templateFiles {
		content, err := templateFS.ReadFile(templateFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read template %s: %w", templateFile, err)
		}

		tmpl, err := template.New(filepath.Base(templateFile)).Parse(string(content))
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", templateFile, err)
		}

		e.templates[templateFile] = tmpl
	}

	return e, nil
}

func (e *Engine) RenderToFile(templateName, outputPath string, data TemplateData) error {
	tmpl, exists := e.templates[templateName]
	if !exists {
		return fmt.Errorf("template %s not found", templateName)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return nil
}
