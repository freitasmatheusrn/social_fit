package database

import (
	"strings"

	"github.com/freitasmatheusrn/social-fit/pkg/rest"
	"github.com/jackc/pgx/v5/pgconn"
)

func GetError(err *pgconn.PgError, columns []string) *rest.ApiErr {
	if err.Code == "23505" {
		for _, c := range columns {
			if strings.Contains(err.ConstraintName, c) {
				cause :=  rest.Causes{
					Field: c,
					Message: "campo já está em uso",
				}
				return rest.NewBadRequestValidationError("campo duplicado", []rest.Causes{cause})
			}
		}
	}
	return rest.NewInternalServerError("erro ao inserir dado")
}
