//go:generate mockgen -source ./repository.go -destination=./mocks/repository.go -package=mock_repository
package repository

import (
	"context"
	"errors"
)

var (
	ErrObjectNotFound = errors.New("object not found")
)

type UsersRepo interface {
	Add(ctx context.Context, user *User) (int64, error)
	GetById(ctx context.Context, id int) (*User, error)
	List(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, user *User) (bool, error)
	Delete(ctx context.Context, id int) (bool, error)
}

type TicketRepo interface {
	Add(ctx context.Context, ticket *Ticket) (int64, error)
	GetById(ctx context.Context, id int) (*Ticket, error)
	GetByUserId(ctx context.Context, UserId int) ([]*Ticket, error)
	List(ctx context.Context) ([]*Ticket, error)
	Update(ctx context.Context, user *Ticket) (bool, error)
	Delete(ctx context.Context, id int) (bool, error)
}
