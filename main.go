package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hisamcode/lenslocked/controllers"
	"github.com/hisamcode/lenslocked/templates"
	"github.com/hisamcode/lenslocked/views"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "layout-page.gohtml", "home.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "layout-page.gohtml", "contact.gohtml"))))

	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.gohtml"))))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	log.Println("Starting server on 3000")
	err := http.ListenAndServe("127.0.0.1:3000", r)
	if err != nil {
		log.Fatal("error listen server :", err)
	}
}
