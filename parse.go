package png

import (
	"context"
	"net"
	"net/url"
	"strconv"
	"time"

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
		updatePort("tcp", u)
		return &HTTPPinger{urlPinger: &urlPinger{url: u}}, nil

	case "tcp":
		fallthrough
	case "tcp4":
		fallthrough
	case "tcp6":
		return parseTCP(u)

	case "mysql":
		if port := u.Port(); port == "" {
			u.Host = u.Hostname() + ":3306"
		}
		updatePort("tcp", u)
		return &MySQLPinger{urlPinger: &urlPinger{url: u}}, nil

	case "postgres":
		if u.RawQuery == "" {
			u.RawQuery = "sslmode=disable"
		}
		if u.Path == "/" {
			u.Path = "/postgres"
		}
		updatePort("tcp", u)
		return &PostgresPinger{urlPinger: &urlPinger{url: u}}, nil

	case "redis":
		return parseRedis(u)

	default:
		return nil, errors.Errorf("invalid scheme: %s", u.Scheme)
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

func parsePort(network, portString string) (int, error) {
	// TODO: invisible context is here!
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	return net.DefaultResolver.LookupPort(ctx, network, portString)
}

func updatePort(network string, u *url.URL) {
	if portString := u.Port(); portString != "" {
		if port, err := parsePort(network, portString); err == nil {
			u.Host = u.Hostname() + ":" + strconv.Itoa(port)
		}
	}
}

func parseTCP(u *url.URL) (Pinger, error) {
	network := u.Scheme
	hostname := u.Hostname()

	portString := u.Port()
	if portString == "" {
		return nil, errors.New("invalid port name: \"\" (empty)")
	}

	port, err := parsePort(network, portString)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid port name: %s", portString)
	}

	return &TCPPinger{network: network, hostname: hostname, port: port}, nil
}

func parseRedis(u *url.URL) (Pinger, error) {
	var password string
	port := 6379 // Redis well-known port
	var db int

	if user := u.User; user != nil {
		password, _ = user.Password()
		// TODO: should we treat username as password? It maybe useful but it maybe break consistent.
	}

	if portString := u.Port(); portString != "" {
		var err error
		port, err = parsePort("tcp", portString)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid port name: %s", portString)
		}
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
		hostname: u.Hostname(),
		port:     port,
		password: password,
		db:       db,
	}, nil
}
