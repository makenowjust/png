package png

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLPinger struct {
	*urlPinger
}

func (p *MySQLPinger) Ping(ctx context.Context) error {
	db, err := sql.Open("mysql", p.urlToDSN())
	if err != nil {
		return errors.Wrap(err, "failed in opening MySQL connection")
	}

	return db.PingContext(ctx)
}

func (p *MySQLPinger) urlToDSN() string {
	u := *p.url
	u.Host = "tcp(" + u.Host + ")"
	return u.String()[len("mysql://"):]
}
