package controllers

import (
	"net/http"

	"github.com/hisamcode/lenslocked/views"
)

type Users struct {
	Templates struct {
		New views.Template
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	// we need view to render
	u.Templates.New.Execute(w, nil)

}
