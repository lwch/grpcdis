package client

import (
	"bufio"
	"fmt"
	"net"

	"github.com/lwch/goredis/code/utils"
	"github.com/lwch/logging"
)

// Client connection client
type Client struct {
	conn  net.Conn
	bufio *bufio.Reader
}

// New new client
func New(conn net.Conn) *Client {
	return &Client{
		conn:  conn,
		bufio: bufio.NewReader(conn),
	}
}

// Close close client
func (cli *Client) Close() error {
	if cli.conn != nil {
		return cli.conn.Close()
	}
	return nil
}

// Run run client
func (cli *Client) Run() {
	defer cli.Close()
	defer utils.Recover("run")
	for {
		argc, err := cli.readArgc()
		if err != nil {
			logging.Error("read argc: %v", err)
			return
		}
		err = cli.readLF()
		if err != nil {
			logging.Error("read crlf: %v", err)
			return
		}
		argv := make([][]byte, argc)
		for i := uint(0); i < argc; i++ {
			argv[i], err = cli.readArgv(i)
			if err != nil {
				logging.Error("read argv(%d): %v", i, err)
				return
			}
		}
		fmt.Println(argv)
	}
}
