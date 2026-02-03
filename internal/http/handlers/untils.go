package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	// Пути к шаблонам (в Docker они будут в /root/web/templates)
	layout := filepath.Join("web/templates/layout.html")
	target := filepath.Join("web/templates", tmpl+".html")

	t, err := template.ParseFiles(layout, target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.ExecuteTemplate(w, "layout", data)
}
