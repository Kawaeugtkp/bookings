package models

import "github.com/Kawaeugtkp/bookings/internal/forms"

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap map[string]int
	FloatMap map[string]float32
	Data map[string]interface{} // interfaceはタイプがはっきりしないものとして
	// 使えるみたいなことを言っています
	CSRFToken string
	Flash string
	Success string
	Warning string
	Error string
	Form *forms.Form
	IsAuthenticated int
}