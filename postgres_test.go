package png

import (
	"testing"

	"bufio"
	"context"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

type PostgresServer struct {
	Close func() error
}

func runPostgresServer(postgresPath, initdbPath, port string) (postgres *PostgresServer, err error) {
	datadir, err := ioutil.TempDir("", "postgres-datadir")
	if err != nil {
		return
	}

	err = exec.Command(initdbPath, "-D", datadir, "-U", "root").Run()
	if err != nil {
		os.RemoveAll(datadir)
		return
	}

	cmd := exec.Command(postgresPath, "-D", datadir, "-p", port)
	r, err := cmd.StderrPipe()
	if err != nil {
		os.RemoveAll(datadir)
		return
	}
	reader := bufio.NewReader(r)

	err = cmd.Start()
	if err != nil {
		os.RemoveAll(datadir)
		r.Close()
		return
	}

	for {
		var line string
		line, err = reader.ReadString('\n')
		if err != nil {
			cmd.Process.Kill()
			r.Close()
			os.RemoveAll(datadir)
			return
		}
		if line == "LOG:  database system is ready to accept connections\n" {
			break
		}
	}

	postgres = &PostgresServer{
		Close: func() error {
			defer os.RemoveAll(datadir)
			defer r.Close()
			return cmd.Process.Kill()
		},
	}
	return
}

func TestPostgresPinger(t *testing.T) {
	lookPath := func(file string) string {
		path, err := exec.LookPath(file)
		if err != nil {
			t.Skip(file + " is not found")
		}
		return path
	}

	postgresPath := lookPath("postgres")
	initdbPath := lookPath("initdb")
	postgres, err := runPostgresServer(postgresPath, initdbPath, "15432")
	if err != nil {
		t.Fatalf("failed in starting postgres: %+#v", err)
	}
	defer postgres.Close()

	t.Run("OK", func(t *testing.T) {
		u, err := url.Parse("postgres://root:@localhost:15432/postgres?sslmode=disable")
		if err != nil {
			panic(err)
		}

		p := &PostgresPinger{url: u}
		err = p.Ping(context.Background())
		if err != nil {
			t.Fatalf("failed in p.Ping(): %+#v", err)
		}
	})

	t.Run("Fail", func(t *testing.T) {
		u, err := url.Parse("postgres://root:@localhost:15432/not_found?sslmode=disable")
		if err != nil {
			panic(err)
		}

		p := &PostgresPinger{url: u}
		err = p.Ping(context.Background())
		if err == nil {
			t.Fatal("succeeded in p.Ping()")
		}

		if msg := err.Error(); !strings.HasPrefix(msg, "failed in Postgres ping: ") {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})
}
