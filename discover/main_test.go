package discover

import (
	"testing"
)

func TestDiscoveryDataExtractors(t *testing.T) {
	var tests = []struct {
		port     string
		fileSize string
	}{
		{"8080", "20"},
		{"6767", "67"},
		{"6789", "897656"},
		{"9090", "99785564"},
	}

	for _, value := range tests {
		var raw = value.port + "\n" + value.fileSize
		var baked = []byte(raw)

		var portBuf []byte
		var fileSizeBuf []byte

		DiscoveryDataExtractor(baked, &portBuf, &fileSizeBuf)
		if string(portBuf) != value.port {
			t.Errorf("port is not eqal to port extracted by DiscoveryDataExtractor. value.port = %s, portBuf = %s", value.port, string(portBuf))
		} else if string(fileSizeBuf) != value.fileSize {
			t.Errorf("fileSize is not eqal to port extracted by DiscoveryDataExtractor. value.fileSize = %s, fileSizeBuf = %s", value.fileSize, string(fileSizeBuf))
		}
	}
}

func TestDiscoveryDataExtractorLowCpuTime(t *testing.T) {
	var tests = []struct {
		port     string
		fileSize string
	}{
		{"8080", "20"},
		{"6767", "67"},
		{"6789", "897656"},
		{"9090", "99785564"},
	}

	for _, value := range tests {
		var raw = value.port + "\n" + value.fileSize
		var baked = []byte(raw)

		var portBuf []byte
		var fileSizeBuf []byte

		DiscoveryDataExtractorLowCpuTime(baked, &portBuf, &fileSizeBuf)
		if string(portBuf) != value.port {
			t.Errorf("port is not eqal to port extracted by DiscoveryDataExtractor. value.port = %s, portBuf = %s", value.port, string(portBuf))
		} else if string(fileSizeBuf) != value.fileSize {
			t.Errorf("fileSize is not eqal to port extracted by DiscoveryDataExtractor. value.fileSize = %s, fileSizeBuf = %s", value.fileSize, string(fileSizeBuf))
		}
	}
}

func BenchmarkDiscoveryDataExtractor(b *testing.B) {

	var raw = "8080\n893893"
	var baked = []byte(raw)

	var portBuf []byte
	var fileSizeBuf []byte
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DiscoveryDataExtractor(baked, &portBuf, &fileSizeBuf)
	}
}
func BenchmarkDiscoveryDataExtractorLowCpuTime(b *testing.B) {
	var raw = "8080\n893893"
	var baked = []byte(raw)

	var portBuf []byte
	var fileSizeBuf []byte
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DiscoveryDataExtractorLowCpuTime(baked, &portBuf, &fileSizeBuf)
	}
}
