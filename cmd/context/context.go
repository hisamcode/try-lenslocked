package main

import (
	stdctx "context"
	"fmt"

	"github.com/hisamcode/lenslocked/context"
	"github.com/hisamcode/lenslocked/models"
)

// type User struct {
// 	name string
// 	age  uint
// }

// type ctxKey string

// const favoriteColorKey ctxKey = "favorite-color"

// func main() {
// 	ctx := context.Background()
// 	user := User{name: "hisam", age: 27}

// 	ctx = context.WithValue(ctx, favoriteColorKey, user)
// 	value := ctx.Value(favoriteColorKey)
// 	// value.(User) is type assertions
// 	fmt.Println(value.(User))
// }

func main() {
	ctx := stdctx.Background()

	user := models.User{
		Email: "hisamcode@gmail.com",
	}

	ctx = context.WithUser(ctx, &user)

	rUser := context.User(ctx)
	fmt.Println(rUser.Email)
}
