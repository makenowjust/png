package png

import (
	"context"
	"database/sql"
	"net/url"

	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLPinger struct {
	url *url.URL
}

func (p *MySQLPinger) Ping(ctx context.Context) error {
	db, _ := sql.Open("mysql", p.urlToDSN())
	// sql.Open() must be succeeded when driver name is correct.
	defer db.Close()

	return errors.Wrap(db.PingContext(ctx), "failed in MySQL ping")
}

func (p *MySQLPinger) urlToDSN() string {
	u := *p.url
	u.Host = "tcp(" + u.Host + ")"
	return u.String()[len("mysql://"):]
}
