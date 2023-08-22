package context

import (
	"context"

	"github.com/hisamcode/lenslocked/models"
)

type key int

const (
	userKey key = iota
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	val := ctx.Value(userKey)
	user, ok := val.(*models.User)
	if !ok {
		// return nil kalo misalkan bukan type models.User
		// the most likely case is that nothing was ever stored in the context,
		// so it doesn't have a type of *models.User. It is also possible that
		// other code in this package wrote an invalid value using the user key.
		return nil
	}
	return user
}
