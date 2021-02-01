package discover

import (
	"net"
	"oDrop/utils"
	"time"
)

func Find() (net.Addr, []byte) {
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
	return addr, buf[:n]
}

func Show(port string) {
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

		_, err = list.Write([]byte(port))
		if err != nil {
			panic(err)
		}
		list.Close()
		time.Sleep(time.Second * 5)
	}
}
