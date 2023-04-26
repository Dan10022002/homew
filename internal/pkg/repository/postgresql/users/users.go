package users

import (
	"context"
	"database/sql"
	"time"

	"homework-week-5/internal/pkg/db"
	"homework-week-5/internal/pkg/repository"
)

type UsersRepo struct {
	db db.DBops
}

func NewUsers(db db.DBops) *UsersRepo {
	return &UsersRepo{db: db}
}

func (r *UsersRepo) Add(ctx context.Context, user *repository.User) (int64, error) {
	var id int64
	err := r.db.ExecQueryRow(ctx, `INSERT INTO users(name, surname, age) VALUES ($1, $2, $3) RETURNING id`,
		user.Name, user.Surname, user.Age).Scan(&id)
	return id, err
}

func (r *UsersRepo) GetById(ctx context.Context, id int) (*repository.User, error) {
	var u repository.User
	err := r.db.Get(ctx, &u, "SELECT id, name, surname, age, created_at, updated_at FROM users WHERE id=$1", id)
	if err == sql.ErrNoRows {
		return nil, repository.ErrObjectNotFound
	}
	return &u, err
}

func (r *UsersRepo) List(ctx context.Context) ([]*repository.User, error) {
	users := make([]*repository.User, 0)
	err := r.db.Select(ctx, &users, "SELECT id, name, surname, age, created_at, updated_at FROM users")
	return users, err
}

func (r *UsersRepo) Update(ctx context.Context, user *repository.User) (bool, error) {
	start := time.Now().UTC()
	result, err := r.db.Exec(ctx,
		"UPDATE users SET name = $1, surname = $2, age = $3, updated_at = $4 WHERE id = $5",
		user.Name, user.Surname, user.Age, start, user.ID)
	return result.RowsAffected() > 0, err
}

func (r *UsersRepo) Delete(ctx context.Context, id int) (bool, error) {
	result, err := r.db.Exec(ctx,
		"DELETE FROM users WHERE id = $1", id)
	return result.RowsAffected() > 0, err
}
