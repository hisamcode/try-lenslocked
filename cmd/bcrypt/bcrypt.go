package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// for i, arg := range os.Args {
	// 	fmt.Println(i, arg)
	// }

	switch os.Args[1] {
	case "hash":
		hash(os.Args[2])
	case "compare":
		compare(os.Args[2], os.Args[3])
	default:
		fmt.Printf("Invalid command :%v\n", os.Args[1])
	}
}

func hash(password string) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error hashing %v\n", password)
		return
	}
	// fmt.Printf("hash the password %q\n", string(hashedBytes))
	fmt.Println(string(hashedBytes))
}

func compare(password, hash string) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Printf("error compare hash and password %v\n", err)
		fmt.Printf("Password is invalid: %v\n", password)
		return
	}
	fmt.Println("Password is correct")
	fmt.Printf("todo compare the password %q with the hash %q\n", password, hash)
}
