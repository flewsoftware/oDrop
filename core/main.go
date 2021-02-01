package core

import (
	"bytes"
	"io"
	"net"
	"oDrop/discover"
	"oDrop/utils"
	"os"
)

// read and broadcast to other receivers (use ReceiveData function for more control)
func Send(callback SendDataCallback, fileLocation string, randomNumber string) error {
	f, err := os.OpenFile(fileLocation, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	s, err := f.Stat()
	dataSize := s.Size()
	err = SendData(callback, f, randomNumber, dataSize)
	return err
}

// receive and write data to a file (use ReceiveData function for more control)
func Receive(writeLocation string, number string, broker dataReceiveBroker) error {
	// receive data
	d, err := ReceiveData(number)
	if err != nil {
		return err
	}

	// create a file
	f, err := os.OpenFile(writeLocation, os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	broker(d, f)
	return nil
}

// A low level receive function used by the higher level Receive
func ReceiveData(number string) (io.Reader, error) {
	add, port := discover.Find()

	c, err := net.Dial("tcp", utils.GetBaseIp(add.String())+":"+string(port))
	if err != nil {
		return nil, err
	}

	_, err = c.Write([]byte(number))
	if err != nil {
		return nil, err
	}

	return c, nil
}

// A low level send function used by the higher level Send
func SendData(callbacks SendDataCallback, reader io.Reader, randomNumber string, dataSize int64) error {
	l, err := net.Listen("tcp", ":6780")
	if err != nil {
		return err
	}
	go func() {
		for {
			c, _ := l.Accept()

			pass := make([]byte, len(randomNumber))

			// extract code from response
			_, err := c.Read(pass)
			if err != nil {
				c.Close()
				break
			} else if bytes.Equal(pass, []byte(randomNumber)) == false {
				c.Close()
				break
			}
			callbacks.DataBroker(c, reader, dataSize)
			c.Close()
			callbacks.SentCallback(c)
		}
	}()

	// sends other listeners the port
	discover.Show("6780")
	return nil
}

// after this function is called net.Conn.Close is called
type dataReceiveBroker func(io.Reader, io.Writer)

type SendDataCallback struct {
	DataBroker   func(net.Conn, io.Reader, int64)
	SentCallback func(net.Conn)
}
