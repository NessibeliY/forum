package template_cache

import (
	"fmt"
	"html/template"
	"path/filepath"

	"01.alem.school/git/nyeltay/forum/internal/models"
)

type TemplateCache map[string]*template.Template

func contains(category string, categories []string) bool {
	for _, c := range categories {
		if c == category {
			return true
		}
	}
	return false
}

func pluck(category *models.Category, categories []*models.Category) []string {
	var names []string
	for _, c := range categories {
		names = append(names, c.Name)
	}
	return names
}

func NewTemplateCache() (map[string]*template.Template, error) {
	pages, err := filepath.Glob("ui/templates/*.page.html")
	if err != nil {
		return nil, fmt.Errorf("load pages: %v", err)
	}

	cache := map[string]*template.Template{}

	// Создаём FuncMap для всех шаблонов
	funcMap := template.FuncMap{
		"contains": contains,
		"pluck":    pluck,
	}
	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(funcMap).ParseFiles(page)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %v", name, err)
		}

		ts, err = ts.ParseGlob("ui/templates/*.partial.html")
		if err != nil {
			return nil, fmt.Errorf("parse %s: %v", name, err)
		}

		cache[name] = ts
	}

	return cache, nil
}
