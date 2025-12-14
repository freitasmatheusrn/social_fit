package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/freitasmatheusrn/social-fit/pkg/bcrypt"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgx.Conn
}

type RepositoryInterface interface{
	Login(ctx context.Context, credentials User) (User, error)
	Signup(ctx context.Context, user User) (User, error)
}

func NewRepo(db *pgx.Conn) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Login(ctx context.Context, credentials User) (User, error) {
	var user User
	err := r.db.QueryRow(
		ctx,
		`SELECT id, name, cpf, email, password, phone, birth_date, created_at
		 FROM users
		 WHERE email = $1`,
		credentials.Email,
	).Scan(&user.ID, &user.Name, &user.Cpf, &user.Email, &user.Password, &user.Phone,  &user.BirthDate, &user.CreatedAt )
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, fmt.Errorf("email não encontrado")
		}
		return User{}, fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	if err := user.ValidatePassword([]byte(credentials.Password)); err != nil {
		return User{}, fmt.Errorf("senha inválida")
	}

	return user, nil
}


func (r *Repository) Signup(ctx context.Context, user User) (User, error) {
	var s User
	hashPassword, err := bcrypt.HashPassword(user.Password)
	if err != nil{
		return User{}, err
	}
	err = r.db.QueryRow(
		ctx,
		`INSERT INTO users (name, email, password, cpf, phone, birthdate)
		 VALUES ($1, $2, $3)
		 RETURNING id name, email, admin`,
		user.Name,
		user.Email,
		hashPassword,
		user.Cpf,
		user.Phone,
		user.BirthDate,
	).Scan(&s.ID, &s.Name, &s.Email, &s.Admin )

	if err != nil {
		return User{}, fmt.Errorf("erro ao inserir usuário: %w", err)
	}

	return s, nil
}