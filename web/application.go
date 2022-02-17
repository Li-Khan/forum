package web

import (
	"html/template"
	"log"

	"github.com/Li-Khan/forum/pkg/models/sqlite"
)

// Application - stores dependencies of the entire web application
type Application struct {
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Forum         *sqlite.ForumModel
	TemplateCache map[string]*template.Template
}
