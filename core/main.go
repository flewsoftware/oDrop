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
	if err != nil {
		return err
	}

	dataSize := s.Size()

	err = SendData(callback, f, randomNumber, dataSize)
	return err
}

// receive and write data to a file (use ReceiveData function for more control)
// (if ip or port is nil it will discover senders in the local network)
func Receive(writeLocation string, number string, broker dataReceiveBroker, ip string, port string, useLowCpuTimeExtractor bool) error {
	// receive data
	d, err, fileSize := ReceiveData(number, ip, port, useLowCpuTimeExtractor)
	if err != nil {
		return err
	}

	// create a file
	f, err := os.OpenFile(writeLocation, os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	broker(d, f, fileSize)
	return nil
}

// A low level receive function used by the higher level Receive
// (if ip or port is nil it will discover senders in the local network)
func ReceiveData(number string, ip string, port string, useLowCpuTimeExtractor bool) (io.Reader, error, []byte) {
	var fileSize = []byte("0")
	if ip == "" || port == "" {
		ipD, portD, fileSizeD := discover.Find(useLowCpuTimeExtractor)
		port = string(portD)
		ip = ipD.String()
		fileSize = fileSizeD
	}

	c, err := net.Dial("tcp", utils.GetBaseIp(ip)+":"+port)
	if err != nil {
		return nil, err, nil
	}

	_, err = c.Write([]byte(number))
	if err != nil {
		return nil, err, nil
	}

	return c, nil, fileSize
}

// A low level send function used by the higher level Send
func SendData(callbacks SendDataCallback, reader io.Reader, randomNumber string, dataSize int64) error {
	l, err := net.Listen("tcp", ":6780")
	if err != nil {
		return err
	}
	StartTcpSever(l, randomNumber, callbacks, reader, dataSize)

	// sends other listeners the port
	discover.Show("6780", dataSize)
	return nil
}

// after this function is called net.Conn.Close is called
type dataReceiveBroker func(io.Reader, io.Writer, []byte)

type SendDataCallback struct {
	DataBroker   func(net.Conn, io.Reader, int64)
	SentCallback func(net.Conn)
}

func StartTcpSever(l net.Listener, randomNumber string, callbacks SendDataCallback, reader io.Reader, dataSize int64) {
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
}
