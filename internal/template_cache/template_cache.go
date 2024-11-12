package template_cache

import (
	"fmt"
	"html/template"
	"path/filepath"
)

type TemplateCache map[string]*template.Template

func NewTemplateCache() (map[string]*template.Template, error) {
	pages, err := filepath.Glob("ui/templates/*.page.html")
	if err != nil {
		return nil, fmt.Errorf("load pages: %v", err)
	}

	cache := map[string]*template.Template{}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %v", name, err)
		}

		// ts, err = ts.ParseGlob("ui/templates/*.layout.html")
		// if err != nil {
		// 	return nil, fmt.Errorf("parse %s: %v", name, err)
		// }

		ts, err = ts.ParseGlob("ui/templates/*.partial.html")
		if err != nil {
			return nil, fmt.Errorf("parse %s: %v", name, err)
		}

		cache[name] = ts
	}

	return cache, nil
}
