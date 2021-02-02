package discover

import (
	"net"
	"oDrop/utils"
	"strconv"
	"strings"
	"time"
)

func Find(useLowCpuTimeExtractor bool) (net.Addr, []byte, []byte) {
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

	trimmedBuf := buf[:n]

	var portBuf []byte
	var fileSizeBuf []byte

	if useLowCpuTimeExtractor {
		DiscoveryDataExtractorLowCpuTime(trimmedBuf, &portBuf, &fileSizeBuf)
	} else {
		DiscoveryDataExtractor(trimmedBuf, &portBuf, &fileSizeBuf)
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

// extracts port and file size from udp discovery data
// this function uses less memory
func DiscoveryDataExtractor(trimmedBuf []byte, portBuf *[]byte, fileSizeBuf *[]byte) {
	var portBufSize = 1
	for i := 0; i < len(trimmedBuf); i++ {
		if trimmedBuf[i] == byte('\n') {
			sBuf := trimmedBuf[portBufSize:]
			for ii := 0; ii < len(sBuf); ii++ {
				*fileSizeBuf = append(*fileSizeBuf, sBuf[ii])
			}
			break
		} else {
			*portBuf = append(*portBuf, trimmedBuf[i])
			portBufSize++
		}
	}
}

// extracts port and file size from udp discovery data
// this function uses less cpu time
func DiscoveryDataExtractorLowCpuTime(trimmedBuf []byte, portBuf *[]byte, fileSizeBuf *[]byte) {
	v := strings.Split(string(trimmedBuf), "\n")

	*portBuf = []byte(v[0])
	*fileSizeBuf = []byte(v[1])
}
