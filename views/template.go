package views

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/hisamcode/lenslocked/context"
	"github.com/hisamcode/lenslocked/models"
)

type public interface {
	Public() string
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl := template.New(patterns[0])
	tpl.Funcs(
		template.FuncMap{
			// i dont know its weird, but this is just for when parsed is not error.
			// the real one is on the bottom, Execute function
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField not implemented")
			},
			"currentUser": func() (template.HTML, error) {
				return "", fmt.Errorf("currentUser not implemented")
			},
			"errors": func() []string {
				return nil
				// return []string{
				// 	"Don't do that!",
				// 	"The email address you provided is already associated with an account.",
				// 	"Something went wrong.",
				// }
			},
		},
	)

	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parse fs template: %w", err)
	}

	return Template{
		htmlTpl: tpl,
	}, nil
}

// func Parse(filepath string) (Template, error) {
// 	tpl, err := template.ParseFiles(filepath)
// 	if err != nil {
// 		return Template{}, fmt.Errorf("parsing template: %w", err)
// 	}

// 	return Template{
// 		htmlTpl: tpl,
// 	}, nil
// }

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {
	// calhound web development in 0118 04:00
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("Cloning template %v", err)
		http.Error(w, "There was an error rendering page", http.StatusInternalServerError)
		return
	}

	errMsgs := errMessages(errs...)
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
			"errors": func() []string {
				return errMsgs
			},
		},
	)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		log.Printf("Executing template: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func errMessages(errs ...error) []string {
	var msgs []string
	for _, err := range errs {
		var pubErr public
		if errors.As(err, &pubErr) {
			msgs = append(msgs, pubErr.Public())
		} else {
			fmt.Println(err)
			msgs = append(msgs, "Something went wrong.")
		}
	}
	return msgs
}
