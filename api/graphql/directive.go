package api

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/loiuscypher/photoview/api/graphql/auth"
)

func IsAdmin(ctx context.Context, obj any, next graphql.Resolver) (res any, err error) {
	user := auth.UserFromContext(ctx)
	if user == nil || user.Admin == false {
		return nil, errors.New("user must be admin")
	}

	return next(ctx)
}

func IsAuthorized(ctx context.Context, obj any, next graphql.Resolver) (res any, err error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, auth.ErrUnauthorized
	}

	return next(ctx)
}
