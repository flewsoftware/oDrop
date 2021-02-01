package discover

import (
	"net"
	"oDrop/utils"
	"strconv"
	"time"
)

func Find() (net.Addr, []byte, []byte) {
	pc, err := net.ListenPacket("udp4", ":8829")
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	buf := make([]byte, 1024)
	n, addr, err := pc.ReadFrom(buf)
	if err != nil {
		panic(err)
	}

	trimedBuffer := buf[:n]

	var portBuf []byte
	var fileSizeBuf []byte
	// mode 0= port /  0 != fileSize
	var mode = 0
	for i := 0; i < len(trimedBuffer); i++ {
		if trimedBuffer[i] == byte('\n') {
			mode = 1
			continue
		}
		if mode == 0 {
			portBuf = append(portBuf, trimedBuffer[i])
		} else {
			fileSizeBuf = append(fileSizeBuf, trimedBuffer[i])
		}
	}

	return addr, portBuf, fileSizeBuf
}

func Show(port string, fileSize int64) {
	local, err := net.ResolveUDPAddr("udp4", ":8829")
	if err != nil {
		panic(err)
	}

	outIP := utils.GetOutboundIP()
	a, _ := utils.LastAddr(net.IPNet{
		IP:   outIP,
		Mask: outIP.DefaultMask(),
	})

	remote, err := net.ResolveUDPAddr("udp4", a.String()+":8829")
	if err != nil {
		panic(err)
	}
	for {
		list, err := net.DialUDP("udp4", local, remote)
		if err != nil {
			panic(err)
		}

		sFileSize := strconv.FormatInt(fileSize, 10)
		_, err = list.Write([]byte(port + "\n" + sFileSize))
		if err != nil {
			panic(err)
		}
		list.Close()
		time.Sleep(time.Second * 5)
	}
}
