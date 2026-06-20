package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"kafe-omborxona/internal/domain"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetAll() ([]domain.User, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, username, password_hash, full_name, role, created_at FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FullName, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepo) GetByID(id int) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(context.Background(),
		`SELECT id, username, password_hash, full_name, role, created_at FROM users WHERE id = $1`, id).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FullName, &u.Role, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByUsername(username string) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(context.Background(),
		`SELECT id, username, password_hash, full_name, role, created_at FROM users WHERE username = $1`, username).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.FullName, &u.Role, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) Create(u *domain.User) error {
	return r.db.QueryRow(context.Background(),
		`INSERT INTO users (username, password_hash, full_name, role) VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		u.Username, u.PasswordHash, u.FullName, u.Role).Scan(&u.ID, &u.CreatedAt)
}

func (r *UserRepo) Update(u *domain.User) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE users SET username = $1, full_name = $2, role = $3, password_hash = CASE WHEN $4 = '' THEN password_hash ELSE $4 END WHERE id = $5`,
		u.Username, u.FullName, u.Role, u.PasswordHash, u.ID)
	return err
}

func (r *UserRepo) Delete(id int) error {
	_, err := r.db.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, id)
	return err
}
