package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
	Age  int
	Bio  string
	Meta UserMeta
}

type UserMeta struct {
	Visits int
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	user := User{
		Name: "John smith",
		Bio:  `<script>alert("haha you");</script>`,
		Age:  111,
		Meta: UserMeta{
			Visits: 10,
		},
	}

	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}

}
