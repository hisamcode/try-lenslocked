package main

import (
	"database/sql"
	"fmt"

	"github.com/hisamcode/lenslocked/models"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgreConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgreConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

func main() {
	cfg := PostgreConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "admin",
		Password: "password",
		Database: "lenslocked",
		SSLMode:  "disable",
	}
	db, err := sql.Open("pgx", cfg.String())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("connected")
	us := models.UserService{
		DB: db,
	}
	user, err := us.Create("bob@bob.com", "bob123")
	if err != nil {
		panic(err)
	}
	fmt.Println(user)

	// create a table
	// _, err = db.Exec(`
	// 	CREATE TABLE IF NOT EXISTS USERS (
	// 		id SERIAL PRIMARY KEY,
	// 		name TEXT,
	// 		email TEXT UNIQUE NOT NULL
	// 	);

	// 	CREATE TABLE IF NOT EXISTS orders (
	// 		id SERIAL PRIMARY KEY,
	// 		user_id INT NOT NULL,
	// 		amount INT,
	// 		description TEXT
	// 	);
	// `)

	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("TABLE USER AND ORDERS CREATED")

	// insert some data
	// name := "bob"
	// email := "bob@gmail.com"
	// row := db.QueryRow(`
	// 	INSERT INTO USERS (name, email)
	// 	VALUES ($1, $2) RETURNING id;
	// `, name, email)

	// var id int
	// if err = row.Scan(&id); err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("User created id = %v\n", id)

	// querying a single record
	// row := db.QueryRow(`
	// 	SELECT name, email FROM users where id = $1
	// `, 10)
	// type User struct {
	// 	Name  string
	// 	Email string
	// }
	// user := User{}
	// if err = row.Scan(&user.Email, &user.Name); err != nil {
	// 	if err == sql.ErrNoRows {
	// 		fmt.Println("error no rows")
	// 	} else {
	// 		panic(err)
	// 	}
	// }

	// fmt.Println(user)

	// userId := 1
	// for i := 1; i <= 5; i++ {
	// 	amount := i * 100
	// 	description := fmt.Sprintf("Fake order #%d", amount)
	// 	_, err := db.Exec(`INSERT INTO orders(user_id, amount, description)
	// 	VALUES ($1, $2, $3)
	// 	`, userId, amount, description)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println("Created fake orders")
	// }

	// type Order struct {
	// 	ID          int
	// 	UserId      int
	// 	Amount      int
	// 	Description string
	// }

	// orders := []Order{}
	// rows, err := db.Query(`
	// 	SELECT id, user_id, amount, description
	// 	FROM orders
	// 	WHERE user_id=$1
	// `, userId)

	// if err != nil {
	// 	panic(err)
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	order := Order{}
	// 	err = rows.Scan(&order.ID, &order.UserId, &order.Amount, &order.Description)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	orders = append(orders, order)
	// }

	// fmt.Println(orders)

}
