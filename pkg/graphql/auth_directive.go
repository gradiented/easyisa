package easyapi

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

var Auth_Directive = func(ctx context.Context, obj interface{}, next graphql.Resolver, role Role) (interface{}, error) {

	return next(ctx)
}
