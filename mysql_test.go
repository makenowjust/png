package png

import (
	"testing"

	"context"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type MySQLServer struct {
	Close func() error
}

func runMySQLServer(mysqldPath, port string) (mysqld *MySQLServer, err error) {
	datadir, err := ioutil.TempDir("", "mysqld-datadir")
	if err != nil {
		return
	}

	err = exec.Command(mysqldPath, "--initialize-insecure", "--datadir="+datadir, "--log-error=mysqld.err").Run()
	if err != nil {
		os.RemoveAll(datadir)
		return
	}

	err = exec.Command(mysqldPath, "--daemonize", "--datadir="+datadir, "--log-error=mysqld.err", "--port="+port, "--pid-file=mysqld.pid").Run()
	if err != nil {
		return
	}

	mysqld = &MySQLServer{
		Close: func() error {
			defer os.RemoveAll(datadir)

			pidFile, err := ioutil.ReadFile(path.Join(datadir, "mysqld.pid"))
			if err != nil {
				return err
			}

			pid, err := strconv.Atoi(strings.TrimSpace(string(pidFile)))
			if err != nil {
				return err
			}

			process, err := os.FindProcess(pid)
			if err != nil {
				return err
			}

			return process.Kill()
		},
	}
	return
}

func TestMySQLPinger(t *testing.T) {
	mysqldPath, err := exec.LookPath("mysqld")
	if err != nil {
		t.Skip("mysqld is not found")
	}

	mysqld, err := runMySQLServer(mysqldPath, "13306")
	if err != nil {
		t.Fatalf("failed in starting mysqld: %+#v", err)
	}
	defer mysqld.Close()

	t.Run("OK", func(t *testing.T) {
		u, err := url.Parse("mysql://root:@localhost:13306/mysql")
		if err != nil {
			panic(err)
		}

		p := &MySQLPinger{url: u}
		err = p.Ping(context.Background())
		if err != nil {
			t.Fatalf("failed in p.Ping(): %+#v", err)
		}
	})

	t.Run("Fail", func(t *testing.T) {
		u, err := url.Parse("mysql://root:@localhost:13306/not_found")
		if err != nil {
			panic(err)
		}

		p := &MySQLPinger{url: u}
		err = p.Ping(context.Background())
		if err == nil {
			t.Fatal("succeeded in p.Ping()")
		}

		if msg := err.Error(); !strings.HasPrefix(msg, "failed in MySQL ping: ") {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	})
}
