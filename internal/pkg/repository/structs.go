package repository

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int          `db:"id"`
	Name      string       `db:"name"`
	Surname   string       `db:"surname"`
	Age       int          `db:"age"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

type Ticket struct {
	ID        int          `db:"id"`
	UserID    int          `db:"user_id"`
	Cost      int          `db:"cost"`
	Place     int          `db:"place"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
