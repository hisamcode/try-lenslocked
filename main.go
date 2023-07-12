package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hisamcode/lenslocked/controllers"
	"github.com/hisamcode/lenslocked/views"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// parse the template
	tpl, err := views.Parse(filepath.Join("templates", "home.gohtml"))
	if err != nil {
		panic(err)
	}
	r.Get("/", controllers.StaticHandler(tpl))

	tpl, err = views.Parse(filepath.Join("templates", "contact.gohtml"))
	if err != nil {
		panic(err)
	}
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl, err = views.Parse(filepath.Join("templates", "faq.gohtml"))
	if err != nil {
		panic(err)
	}
	r.Get("/faq", controllers.StaticHandler(tpl))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	log.Println("Starting server on 3000")
	err = http.ListenAndServe("127.0.0.1:3000", r)
	if err != nil {
		log.Fatal("error listen server :", err)
	}
}
