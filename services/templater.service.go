package services

import (
	"bytes"
	"html/template"
	"path"

	"go.uber.org/zap"
)

type Templater interface {
	Template(filePath string, data interface{}) (string, error)
}

type TemplateService struct {
	logger     *zap.SugaredLogger
	assetsPath string
}

func NewTemplaterService(assetsPath string, logger *zap.SugaredLogger) Templater {
	return &TemplateService{
		logger:     logger,
		assetsPath: assetsPath,
	}
}

func (t *TemplateService) Template(filePath string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(path.Join(t.assetsPath, filePath))
	if err != nil {
		t.logger.Errorw("Failed to parse template", "error", err)
		return "", err
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		t.logger.Errorw("Failed to execute template", "error", err)
		return "", err
	}

	return tpl.String(), nil
}
