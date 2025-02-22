package migrate

import (
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
)

func Migrate(dir string) error {
	db, err := migrate.New(dir, "mysql://" + os.Getenv("MYSQL_DSN"))

	if err != nil {
		return errors.Wrap(err, "migrate.New")
	}

	err = db.Up()

	if err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "db.Up")
	}

	serr, dberr := db.Close()

	if serr != nil {
		return errors.Wrap(serr, "db.Close source error")
	}

	if dberr != nil {
		return errors.Wrap(dberr, "db.Close database error")
	}

	return nil
}
