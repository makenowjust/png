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
	db, _ := sql.Open("postgres", p.url.String())
	// sql.Open() must be succeeded when driver name is correct.
	defer db.Close()

	return errors.Wrap(db.PingContext(ctx), "failed in Postgres ping")
}
