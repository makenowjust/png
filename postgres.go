package png

import (
	"context"
	"database/sql"
	"net/url"

	"github.com/pkg/errors"

	_ "github.com/lib/pq"
)

type PostgresPinger struct {
	url *url.URL
}

func (p *PostgresPinger) Ping(ctx context.Context) error {
	db, err := sql.Open("postgres", p.url.String())
	if err != nil {
		return errors.Wrap(err, "failed in opening Postgres connection")
	}

	return errors.Wrap(db.PingContext(ctx), "failed in Postgres ping")
}
