package utils

import (
	"testing"
)

func TestGetBaseIp(t *testing.T) {
	var tests = []struct {
		addrWithPort string
		ipAddr       string
	}{
		{"localhost:8080", "localhost"},
		{"192.168.8.1:8989", "192.168.8.1"},
		{"192.168.8.200:8076", "192.168.8.200"},
	}

	for _, value := range tests {
		if a := GetBaseIp(value.addrWithPort); a != value.ipAddr {
			t.Errorf("ipaddress not eqal to GetBaseIp result: %s != %s", a, value.ipAddr)
		}
	}
}

func BenchmarkGetBaseIp(b *testing.B) {
	var p = "localhost:8080"

	for i := 0; i < b.N; i++ {
		_ = GetBaseIp(p)
	}

}
