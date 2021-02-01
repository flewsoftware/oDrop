package utils

import (
	"encoding/binary"
	"errors"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func DoesFileExist(fileName string) bool {
	i, err := os.Stat(fileName)
	if err != nil {
		return false
	} else if i.IsDir() {
		return false
	}
	return true
}

func ModeToSimple(mode string) string {
	if mode == "send" || mode == "s" {
		return "s"
	} else {
		return "r"
	}
}

func GetBaseIp(addr string) string {
	return strings.Split(addr, ":")[0]
}

func LastAddr(n net.IPNet) (net.IP, error) {
	if n.IP.To4() == nil {
		return net.IP{}, errors.New("does not support IPv6 addresses.")
	}
	ip := make(net.IP, len(n.IP.To4()))
	binary.BigEndian.PutUint32(ip, binary.BigEndian.Uint32(n.IP.To4())|^binary.BigEndian.Uint32(net.IP(n.Mask).To4()))
	return ip, nil
}

func GetOutboundIP() net.IP {

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func RemoveWhitespace(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

func GetRandomNumber() string {
	rand.Seed(time.Now().UnixNano())
	rand.Int()
	r := strings.Split(strconv.Itoa(rand.Int()), "")
	var sRandomNumber string

	for i := 0; i <= 6; i++ {
		sRandomNumber += r[i]
	}
	return sRandomNumber
}
