// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"context"
)

type Querier interface {
	GetPublicKey(ctx context.Context, username *string) (string, error)
	RegisterUser(ctx context.Context, arg RegisterUserParams) error
}

var _ Querier = (*Queries)(nil)
