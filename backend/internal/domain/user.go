package domain

import "time"

type UserRole string

const (
	RoleAdmin            UserRole = "admin"
	RoleWarehouseManager UserRole = "warehouse_manager"
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	FullName string   `json:"full_name"`
	Role     UserRole `json:"role"`
}

type UpdateUserRequest struct {
	Username string   `json:"username"`
	FullName string   `json:"full_name"`
	Role     UserRole `json:"role"`
	Password string   `json:"password,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UserRepository interface {
	GetAll() ([]User, error)
	GetByID(id int) (*User, error)
	GetByUsername(username string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id int) error
}
