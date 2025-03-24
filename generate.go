package sn

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

type style string

type length int

const (
	Date    style = "date"
	Week    style = "week"
	Quarter style = "quarter"
)

const (
	Length24 length = 24
	Length22 length = 22
)

type sn struct {
	time   time.Time
	style  style
	length length
}

type Option func(*sn)

type generate interface {
	generate(time.Time) string
}

func WithStyle(style style) Option {
	return func(s *sn) {
		s.style = style
	}
}

func WithLength(length length) Option {
	return func(s *sn) {
		s.length = length
	}
}

func WithTime(time time.Time) Option {
	return func(s *sn) {
		s.time = time
	}
}

var num int64

// Generate 生成24位订单号 [前面17位代表时间精确到毫秒，中间3位代表进程id，最后4位代表序号]
func Generate(opts ...Option) string {
	s := &sn{
		time:   time.Now(),
		style:  Date,
		length: Length24,
	}
	for _, opt := range opts {
		opt(s)
	}

	str := factory[s.style].generate(s.time)
	ms := (fmt.Sprintf("%d", s.time.UnixMilli()))[10:]
	p := fmt.Sprintf("%03d", os.Getpid()%1000)
	n := fmt.Sprintf("%04d", atomic.AddInt64(&num, 1)%10000)

	str = fmt.Sprintf("%s%s%s%s", str, ms, p, n)

	if s.length == Length22 {
		str = str[2:]
	}

	return str
}

type date struct{}

type week struct{}

type quarter struct{}

func (*date) generate(t time.Time) string {
	return t.Format("20060102150405")
}

func (*week) generate(t time.Time) string {
	year, _week := t.ISOWeek()
	return fmt.Sprintf("%d%02d%s", year, _week, t.Format("02150405"))
}

func (*quarter) generate(t time.Time) string {
	month := t.Month()
	var q int
	var m int
	if month >= 1 && month <= 3 {
		q = 1
		m = int(month)
	} else if month >= 4 && month <= 6 {
		q = 2
		m = int(month) - 3
	} else if month >= 7 && month <= 9 {
		q = 3
		m = int(month) - 6
	} else {
		q = 4
		m = int(month) - 9
	}
	return fmt.Sprintf("%d%d%d%s", t.Year(), q, m, t.Format("02150405"))
}

var factory map[style]generate

func init() {
	factory = map[style]generate{
		Date:    &date{},
		Week:    &week{},
		Quarter: &quarter{},
	}
}
