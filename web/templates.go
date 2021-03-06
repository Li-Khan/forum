package web

import (
	"html/template"
	"path/filepath"

	"github.com/Li-Khan/forum/pkg/models"
)

// templateData ...
type templateData struct {
	User       models.User
	Post       models.Post
	Posts      []models.Post
	Errors     models.Errors
	NumOfPosts int
	Tags       []string
	IsSession  bool
	Error      string
}

// NewTemplateCache - initializes the template cache
func NewTemplateCache(dir string) (map[string]*template.Template, error) {
	// Инициализирую новую карту, которая будет хранить кэш.
	cache := map[string]*template.Template{}

	// filepath.Glob - получает срез всех файловых путей с расширением "page.html".
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Использую метод ParseGlob для добавления всех каркасных шаблонов.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}

		// Использую метод ParseGlob для добавления всех вспомогательных шаблонов.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}

		// Добавляю полученный набор шаблонов в кэш.
		cache[name] = ts
	}
	return cache, nil
}
