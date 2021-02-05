package speedwrap

import (
	"math"
	"time"
)

type Writer interface {
	Write(p []byte) (n int, err error)
}

type Reader interface {
	Read(p []byte) (n int, err error)
}

type SW struct {
	bytes     int64
	startTime time.Time
}

func (s *SW) SetStartTime() {
	s.startTime = time.Now()
}

func (s *SW) Read(b []byte) (n int, err error) {
	n = len(b)
	s.bytes += int64(n)
	return
}
func (s *SW) Write(b []byte) (n int, err error) {
	n = len(b)
	s.bytes += int64(n)
	return
}

func (s *SW) GetSpeed() float64 {
	return float64(s.bytes) / time.Since(s.startTime).Seconds()
}
func (s *SW) GetSpeedRound() int64 {
	return int64(math.Round(s.GetSpeed()))
}
