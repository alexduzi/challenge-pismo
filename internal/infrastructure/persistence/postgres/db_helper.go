package postgres

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alexduzi/challengepismo/internal/infrastructure/exception"
	"github.com/jackc/pgx/v5/pgconn"
)

func handlePgError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			if strings.Contains(pgErr.ConstraintName, "document_number") {
				return exception.ErrDuplicateDocument
			}
			if strings.Contains(pgErr.ConstraintName, "email") {
				return exception.ErrDuplicateEmail
			}
		}
	}
	return fmt.Errorf("%w: %v", exception.ErrDatabaseError, err)
}
