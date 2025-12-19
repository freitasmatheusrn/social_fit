package user

import (
	"context"
	"errors"
	"log"

	"github.com/freitasmatheusrn/social-fit/internal/database"
	"github.com/freitasmatheusrn/social-fit/pkg/bcrypt"
	"github.com/freitasmatheusrn/social-fit/pkg/rest"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	db *pgx.Conn
}

type RepositoryInterface interface {
	Login(ctx context.Context, credentials User) (User, rest.ApiErr)
	Signup(ctx context.Context, user User) (User, rest.ApiErr)
}

func NewRepo(db *pgx.Conn) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Login(ctx context.Context, credentials User) (User, *rest.ApiErr) {
	var user User
	err := r.db.QueryRow(
		ctx,
		`SELECT id, name, email, admin, password
		 FROM users
		 WHERE email = $1`,
		credentials.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Admin, &user.Password )
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			cause := rest.Causes{
				Field:   "email",
				Message: "email não encontrado",
			}
			return User{}, rest.NewBadRequestValidationError("email não encontrado", []rest.Causes{cause})
		}
		log.Println(err)

		return User{}, rest.NewInternalServerError("erro ao buscar usuário")
	}

	if err := user.ComparePassword([]byte(credentials.Password)); err != nil {
		cause := rest.Causes{
			Field:   "password",
			Message: "senha incorreta",
		}
		return User{}, rest.NewBadRequestValidationError("senha incorreta", []rest.Causes{cause})
	}

	return user, nil
}

func (r *Repository) Signup(ctx context.Context, user User) (User, *rest.ApiErr) {
	var u User
	hashPassword, err := bcrypt.HashPassword(user.Password)
	if err != nil {
		return User{}, rest.NewInternalServerError("erro do servidor")
	}
	err = r.db.QueryRow(
		ctx,
		`INSERT INTO users (name, email, password, cpf, phone, birth_date)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, name, email, admin`,
		user.Name,
		user.Email,
		hashPassword,
		user.Cpf,
		user.Phone,
		user.BirthDate,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Admin)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return User{}, database.GetError(pgErr, []string{"email"})
		}
		return User{}, rest.NewInternalServerError("erro interno")
	}

	return u, nil
}
