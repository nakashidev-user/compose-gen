package infrastructure

import (
	"bytes"
	"context"
	"embed"
	"text/template"

	"compose-gen/internal/domain"
)

//go:embed templates/docker-compose.yml.tmpl
var templateFS embed.FS

type TemplateRepositoryImpl struct{}

func NewTemplateRepository() *TemplateRepositoryImpl {
	return &TemplateRepositoryImpl{}
}

func (r *TemplateRepositoryImpl) Generate(ctx context.Context, config domain.ComposeConfig) (string, error) {
	tmplContent, err := templateFS.ReadFile("templates/docker-compose.yml.tmpl")
	if err != nil {
		return "", domain.ErrTemplateNotFound
	}

	tmpl, err := template.New("docker-compose").Parse(string(tmplContent))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return "", err
	}

	return buf.String(), nil
}