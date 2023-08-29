package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/hisamcode/lenslocked/controllers"
	"github.com/hisamcode/lenslocked/migrations"
	"github.com/hisamcode/lenslocked/models"
	"github.com/hisamcode/lenslocked/templates"
	"github.com/hisamcode/lenslocked/views"
	"github.com/joho/godotenv"
)

type config struct {
	PSQL models.PostgreConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (*config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	// PSQL
	cfg.PSQL = models.DefaultPostgresConfig()

	//SMTP
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	// CSRF
	cfg.CSRF.Key = "32-byte-long-auth-key"
	cfg.CSRF.Secure = false

	// Server
	cfg.Server.Address = "127.0.0.1:3000"

	return &cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	// Setup the database
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Setup services
	userService := &models.UserService{DB: db}
	sessionService := &models.SessionService{DB: db}
	pwResetService := &models.PasswordResetService{DB: db}
	emailService := models.NewEmailService(cfg.SMTP)
	galleryService := &models.GalleryService{DB: db}

	// Setup middleware
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		// fix this before deploy
		csrf.Secure(cfg.CSRF.Secure),
		csrf.Path("/"),
	)

	// Setup controllers
	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
	}

	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"signup.gohtml", "tailwind.gohtml",
	))
	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"signin.gohtml", "tailwind.gohtml",
	))

	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(
		templates.FS,
		"forgot-pw.gohtml", "tailwind.gohtml",
	))

	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(
		templates.FS,
		"check-your-email.gohtml", "tailwind.gohtml",
	))

	usersC.Templates.ResetPassword = views.Must(views.ParseFS(
		templates.FS,
		"reset-pw.gohtml", "tailwind.gohtml",
	))

	galleriesC := controllers.Galleries{
		GalleryService: galleryService,
	}
	galleriesC.Template.New = views.Must(views.ParseFS(
		templates.FS,
		"galleries/new.gohtml", "tailwind.gohtml",
	))
	galleriesC.Template.Edit = views.Must(views.ParseFS(
		templates.FS,
		"galleries/edit.gohtml", "tailwind.gohtml",
	))
	galleriesC.Template.Index = views.Must(views.ParseFS(
		templates.FS,
		"galleries/index.gohtml", "tailwind.gohtml",
	))
	galleriesC.Template.Show = views.Must(views.ParseFS(
		templates.FS,
		"galleries/show.gohtml", "tailwind.gohtml",
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
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)
	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)
	// r.Get("/users/me", usersC.CurrentUser)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
		r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "hello")
		})
	})
	// r.Get("/galleries/new", galleriesC.New)
	r.Route("/galleries", func(r chi.Router) {
		r.Get("/{id}", galleriesC.Show)
		r.Get("/{id}/images/{filename}", galleriesC.Image)
		r.Group(func(r chi.Router) {
			r.Use(umw.RequireUser)
			r.Get("/", galleriesC.Index)
			r.Get("/new", galleriesC.New)
			r.Post("/", galleriesC.Create)
			r.Get("/{id}/edit", galleriesC.Edit)
			r.Post("/{id}", galleriesC.Update)
			r.Post("/{id}/delete", galleriesC.Delete)
			r.Post("/{id}/images", galleriesC.UploadImage)
			r.Post("/{id}/images/{filename}/delete", galleriesC.DeleteImage)
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// Start the server
	log.Printf("Starting server on %s...\n", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}
}

// func TimerMiddleware(h http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
// 		h(w, r)
// 		fmt.Println("Request time: ", time.Since(start))
// 	}
// }
