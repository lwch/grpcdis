package client

import (
	"io"
	"strconv"
)

func (cli *Client) readArgc() (uint, error) {
	str, err := cli.bufio.ReadString('\r')
	if err != nil {
		return 0, err
	}
	if str[0] != '*' {
		return 0, errCount
	}
	n, err := strconv.ParseUint(str[1:len(str)-1], 10, 64)
	if err != nil {
		return 0, err
	}
	err = cli.readLF()
	if err != nil {
		return 0, err
	}
	return uint(n), nil
}

func (cli *Client) readArgv(i uint) ([]byte, error) {
	str, err := cli.bufio.ReadString('\r')
	if err != nil {
		return nil, err
	}
	if str[0] != '$' {
		return nil, errSize
	}
	size, err := strconv.ParseUint(str[1:len(str)-1], 10, 64)
	if err != nil {
		return nil, err
	}
	err = cli.readLF()
	if err != nil {
		return nil, err
	}
	data := make([]byte, size)
	_, err = io.ReadFull(cli.bufio, data)
	if err != nil {
		return nil, err
	}
	err = cli.readCRLF()
	if err != nil {
		return nil, err
	}
	return data, nil
}
