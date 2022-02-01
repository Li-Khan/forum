package web

import (
	"log"

	"github.com/Li-Khan/forum/pkg/models/sqlite"
)

// Application ...
type Application struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Forum    *sqlite.ForumModel
}
