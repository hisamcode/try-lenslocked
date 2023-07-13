package controllers

import (
	"html/template"
	"net/http"

	"github.com/hisamcode/lenslocked/views"
)

func StaticHandler(tpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}

func FAQ(tpl views.Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "Is there a free version ?",
			Answer:   "Yes! we offer a free trial for 30 days on any paid plans.",
		},
		{
			Question: "What are your support hours?",
			Answer:   "We have support staff answering emails 24/7, though response time may be a bit slower on weekends.",
		},
		{
			Question: "How do i contact support?",
			Answer:   `Email us - <a href="mailto:hisamcode@gmail.com">hisamcode@gmail.com</a>`,
		},
		{
			Question: "Do you have office located?",
			Answer:   `Our entire team is remotely`,
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, questions)
	}
}
