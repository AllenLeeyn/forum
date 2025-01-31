package handlers

import (
	"net/http"
	"text/template"
)

func aboutHandler(w http.ResponseWriter, r *http.Request) {

	// Parse template
	tmpl, err := template.ParseFiles("templates/about.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute template
	err = tmpl.Execute(w, "about.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
