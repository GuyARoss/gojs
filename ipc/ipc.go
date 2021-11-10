package ipc

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
)

type IPCClient struct {
	path string
	Name string
}

type IncomingData struct {
	OpName string
	Data   string
}

type IPCListener struct {
	IncomingData chan IncomingData
}

func parseIncomingBytes(data []byte) *IncomingData {
	f := strings.SplitN(string(data), ":", 2)

	return &IncomingData{
		OpName: f[0],
		Data:   f[1],
	}
}

var ErrNoClient = errors.New("client is nil")

func (c *IPCClient) Write(operationName string, data string) error {
	if c == nil {
		return ErrNoClient
	}

	f, err := os.OpenFile(c.path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return err
	}

	// note to future dev, DO NOT remove the \n.. this caused hours of headache..
	f.WriteString(fmt.Sprintf("%s:%s\n", operationName, data))

	return nil
}

func (c *IPCClient) blockingListener(handler func(incoming *IncomingData)) {
	file, err := os.OpenFile(c.path, os.O_CREATE, os.ModeNamedPipe)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadBytes('\n')
		if err == nil {
			incoming := parseIncomingBytes(line)

			handler(incoming)
		}
	}
}

func (c *IPCClient) Listen() (*IPCListener, error) {
	incomingChan := make(chan IncomingData)

	go func(ic chan IncomingData) {
		c.blockingListener(func(data *IncomingData) {
			ic <- *data
		})
	}(incomingChan)

	return &IPCListener{
		IncomingData: incomingChan,
	}, nil
}

func (c *IPCClient) Dispose() {
	os.Remove(c.path)
}

func New(name string) (*IPCClient, error) {
	path := fmt.Sprintf("/tmp/%s", name)

	_, err := os.Stat(path)
	if err != nil {
		os.Remove(path)
		fifoerr := syscall.Mkfifo(path, 0666)

		if fifoerr != nil {
			return nil, fifoerr
		}
	}

	return &IPCClient{
		path: path,
		Name: name,
	}, nil
}
