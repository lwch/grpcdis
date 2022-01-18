package client

import (
	"encoding/hex"
	"fmt"
	"io"
	"strconv"

	"github.com/lwch/goredis/code/utils"
)

func (cli *Client) readArgc() (uint, error) {
	str, err := cli.bufio.ReadBytes('\r')
	if err != nil {
		return 0, err
	}
	fmt.Println(hex.Dump(str))
	return 0, nil
	// op, err := cli.bufio.ReadByte()
	// if err != nil {
	// 	return 0, err
	// }
	// if op != '*' {
	// 	return 0, errCount
	// }
	// var data []byte
	// for {
	// 	ch, err := cli.bufio.ReadByte()
	// 	if err != nil {
	// 		return 0, err
	// 	}
	// 	if !utils.IsNumber(ch) {
	// 		cli.bufio.Discard()
	// 		break
	// 	}
	// 	data = append(data, ch)
	// }
	// n, err := strconv.ParseUint(string(data), 10, 64)
	// if err != nil {
	// 	return 0, err
	// }
	// return uint(n), nil
}

func (cli *Client) readArgv(i uint) ([]byte, error) {
	op, err := cli.bufio.ReadByte()
	if err != nil {
		return nil, err
	}
	if op != '$' {
		return nil, errSize
	}
	var data []byte
	for {
		ch, err := cli.bufio.ReadByte()
		if err != nil {
			return nil, err
		}
		if !utils.IsNumber(ch) {
			break
		}
		data = append(data, ch)
	}
	size, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		return nil, err
	}
	err = cli.readLF()
	if err != nil {
		return nil, err
	}
	data = make([]byte, size)
	_, err = io.ReadFull(cli.bufio, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
