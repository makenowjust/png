package png

import (
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

func Parse(rawurl string) (Pinger, error) {
	if rawurl == "" {
		return nil, errors.New("invalid URL: \"\" (empty)")
	}

	u, err := parseURL(rawurl)
	if err != nil {
		return nil, err
	}

	// When hostname is not specified, sets `127.0.0.1`.
	if u.Hostname() == "" {
		host := "127.0.0.1"
		if port := u.Port(); port != "" {
			host += ":" + port
		}
		u.Host = host
	}

	switch u.Scheme {
	case "http":
		fallthrough
	case "https":
		return &HTTPPinger{url: u}, nil
	case "ws":
		fallthrough
	case "wss":
		return &WebSocketPinger{url: u}, nil

	case "tcp":
		fallthrough
	case "tcp4":
		fallthrough
	case "tcp6":
		return &TCPPinger{network: u.Scheme, addr: u.Host}, nil

	case "mysql":
		if port := u.Port(); port == "" {
			u.Host = u.Hostname() + ":3306"
		}
		return &MySQLPinger{url: u}, nil

	case "postgres":
		if u.RawQuery == "" {
			u.RawQuery = "sslmode=disable"
		}
		if u.Path == "/" {
			u.Path = "/postgres"
		}
		return &PostgresPinger{url: u}, nil

	case "redis":
		return parseRedis(u)

	case "amqp":
		return &AMQPPinger{url: u}, nil

	default:
		return nil, errors.Errorf("unknown scheme: %s", u.Scheme)
	}
}

func parseURL(rawurl string) (u *url.URL, err error) {
	u, err = url.Parse(rawurl)
	if err != nil {
		err = errors.Wrapf(err, "failed in parsing URL: %#v", rawurl)
	} else {
		// - `u.Opaque != ""` is for `localhost:8080` case.
		// - `u.Scheme != ""` is for `localhost` case.
		if u.Opaque != "" || u.Scheme == "" {
			rawurl = "http://" + rawurl

			u, err = url.Parse(rawurl)
			err = errors.Wrapf(err, "failed in parsing URL: %#v", rawurl)
		}
	}

	return
}

func parseRedis(u *url.URL) (Pinger, error) {
	var password string
	var db int

	if user := u.User; user != nil {
		password, _ = user.Password()
		// TODO: should we treat username as password? It maybe useful but it maybe break consistent.
	}

	if path := u.Path; len(path) >= 2 {
		path = path[1:] // skip `/`

		var err error
		db, err = strconv.Atoi(path)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid db number: %#v", path)
		}
	}

	return &RedisPinger{
		addr:     u.Host,
		password: password,
		db:       db,
	}, nil
}
