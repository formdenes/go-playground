package dberror

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrUniqueViolation     = errors.New("unique violation")
	ErrNoRows              = errors.New("no rows")
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrException           = errors.New("exception raised")
)

func DbErrorFromPq(err error) error {
	log.Printf("Db error:  %v", err)
	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, sql.ErrNoRows) {
		log.Print("not not found")
		return fmt.Errorf("%w%w", ErrNoRows, err)
	}
	var pgerr *pgconn.PgError
	ok := errors.As(err, &pgerr)
	if !ok {
		log.Print("not postgres error")
		return err
	}
	switch pgerr.Code {
	case "23503":
		return fmt.Errorf("%w%w", ErrForeignKeyViolation, err)
	case "23505":
		return fmt.Errorf("%w%w", ErrUniqueViolation, err)
	case "P0001":
		return fmt.Errorf("%w%w", ErrUniqueViolation, err)
	}

	log.Printf("pqErr: %v", pgerr)
	log.Print("not handled postgres error")
	return err
}
