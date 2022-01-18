package client

func (cli *Client) readLF() error {
	lf, err := cli.bufio.ReadByte()
	if err != nil {
		return err
	}
	if lf != '\n' {
		return errLF
	}
	return nil
}
