package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/hisamcode/lenslocked/controllers"
	"github.com/hisamcode/lenslocked/migrations"
	"github.com/hisamcode/lenslocked/models"
	"github.com/hisamcode/lenslocked/templates"
	"github.com/hisamcode/lenslocked/views"
)

func main() {

	// Setup the database
	cfg := models.DefaultPostgresConfig()
	fmt.Println(cfg)
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Setup services
	userService := models.UserService{DB: db}
	sessionService := models.SessionService{DB: db}

	// Setup middleware
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	csrfKey := "32-byte-long-auth-key"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		// fix this before deploy
		csrf.Secure(false),
	)

	// Setup controllers
	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}

	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"signup.gohtml", "tailwind.gohtml",
	))
	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"signin.gohtml", "tailwind.gohtml",
	))

	// Setup our router and routes
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(csrfMw)
	r.Use(umw.SetUser)

	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(
			templates.FS,
			"home.gohtml", "tailwind.gohtml",
		))))

	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(
			templates.FS,
			"contact.gohtml", "tailwind.gohtml",
		))))

	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(
			templates.FS,
			"faq.gohtml",
			"tailwind.gohtml",
		))))

	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProceessSignOut)
	// r.Get("/users/me", usersC.CurrentUser)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
		r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "hello")
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// Start the server
	log.Println("Starting server on 3000")
	err = http.ListenAndServe("127.0.0.1:3000", r)
	if err != nil {
		log.Fatal("error listen server :", err)
	}
}

// func TimerMiddleware(h http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
// 		h(w, r)
// 		fmt.Println("Request time: ", time.Since(start))
// 	}
// }
