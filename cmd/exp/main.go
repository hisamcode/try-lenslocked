package main

import (
	"errors"
	"fmt"

	"github.com/hisamcode/lenslocked/models"
)

func main() {
	p := models.PostgreConfig{Host: "haha"}
	fmt.Printf("%s", p)
	// err := CreateOrg()
	// if err != nil {
	// 	err = errors.Unwrap(err)
	// 	fmt.Println(err)
	// }
}

func CreateOrg() error {
	err := CreateUser()
	if err != nil {
		return fmt.Errorf("create org: %w", err)
	}
	return nil
}

func Connect() error {
	return errors.New("connection failed")
}

func CreateUser() error {
	err := Connect()
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}
