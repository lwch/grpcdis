package client

import (
	"bufio"
	"net"

	"github.com/lwch/goredis/code/command"
	"github.com/lwch/goredis/code/command/server"
	"github.com/lwch/goredis/code/utils"
	"github.com/lwch/logging"
)

// Client connection client
type Client struct {
	conn   net.Conn
	bufio  *bufio.Reader
	writer *command.LockWriter
	cmds   *server.Command
}

// New new client
func New(conn net.Conn, cmds *server.Command) *Client {
	return &Client{
		conn:   conn,
		bufio:  bufio.NewReader(conn),
		writer: command.NewWriter(conn),
		cmds:   cmds,
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
		argv := make([][]byte, argc)
		for i := uint(0); i < argc; i++ {
			argv[i], err = cli.readArgv(i)
			if err != nil {
				logging.Error("read argv(%d): %v", i, err)
				return
			}
		}
		cmd := cli.cmds.Lookup(string(argv[0]))
		if cmd == nil {
			// TODO: reply error
			logging.Error("command [%s] not supported", string(argv[0]))
			return
		}
		err = cmd.Run(argv[1:], cli.writer)
		if err != nil {
			// TODO: reply error
			logging.Error("run command [%s]: %v", string(argv[0]), err)
			return
		}
	}
}
