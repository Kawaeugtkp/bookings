package config

import (
	"html/template"
	"log"

	"github.com/Kawaeugtkp/bookings/internal/models"
	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the application config
type AppConfig struct {
	Usercache bool
	TemplateCache map[string]*template.Template
	InfoLog *log.Logger
	ErrorLog *log.Logger
	InProduction bool
	Session *scs.SessionManager
	MailChan chan models.MailData
}