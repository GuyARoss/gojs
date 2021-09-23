package ipc

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
)

type IPCClient struct {
	path string
	Name string
	file *os.File
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

func (c *IPCClient) Write(operationName string, data string) error {
	// note to future dev, DO NOT remove the \n.. this caused hours of headache..
	c.file.WriteString(fmt.Sprintf("%s:%s\n", operationName, data))

	return nil
}

func (c *IPCClient) blockingListener(handler func(incoming *IncomingData)) {
	reader := bufio.NewReader(c.file)

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

	handleIncomingChan := func(data *IncomingData) {
		incomingChan <- *data
	}

	go func() {
		c.blockingListener(handleIncomingChan)
	}()

	return &IPCListener{
		IncomingData: incomingChan,
	}, nil
}

func New(name string) (*IPCClient, error) {
	path := fmt.Sprintf("/tmp/%s", name)

	_, err := os.Stat(path)
	if err == nil {
		os.Remove(path)
	}

	fifoerr := syscall.Mkfifo(path, 0666)
	if fifoerr != nil {
		return nil, fifoerr
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	return &IPCClient{
		file: file,
		path: path,
		Name: name,
	}, nil
}
