package png

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	_ "github.com/lib/pq"
)

type PostgresPinger struct {
	*urlPinger
}

func (p *PostgresPinger) Ping(ctx context.Context) error {
	db, err := sql.Open("postgres", p.url.String())
	if err != nil {
		return errors.Wrap(err, "failed in opening Postgres connection")
	}

	return db.PingContext(ctx)
}
