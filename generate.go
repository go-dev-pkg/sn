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
	Date style = "date"
	Week style = "week"
)

const (
	Length24 length = 24
	Length22 length = 22
)

type sn struct {
	style  style
	length length
}

type Option func(*sn)

type generate interface {
	generate() string
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

var num int64

func Generate(opts ...Option) string {
	s := &sn{
		style:  Date,
		length: Length24,
	}
	for _, opt := range opts {
		opt(s)
	}

	str := factory[s.style].generate()

	if s.length == Length22 {
		str = str[2:]
	}

	return str
}

type date struct{}

type week struct{}

func (*date) generate() string {
	t := time.Now()

	d := t.Format("20060102150405")
	ms := (fmt.Sprintf("%d", t.UnixMilli()))[10:]
	p := fmt.Sprintf("%03d", os.Getpid()%1000)
	n := fmt.Sprintf("%04d", atomic.AddInt64(&num, 1)%10000)

	return fmt.Sprintf("%s%s%s%s", d, ms, p, n)
}

func (*week) generate() string {
	t := time.Now()

	year, _week := t.ISOWeek()
	w := fmt.Sprintf("%02d", _week)
	d := t.Format("02150405")
	ms := (fmt.Sprintf("%d", t.UnixMilli()))[10:]
	p := fmt.Sprintf("%03d", os.Getpid()%1000)
	n := fmt.Sprintf("%04d", atomic.AddInt64(&num, 1)%10000)

	return fmt.Sprintf("%d%s%s%s%s%s", year, w, d, ms, p, n)
}

var factory map[style]generate

func init() {
	factory = map[style]generate{
		Date: &date{},
		Week: &week{},
	}
}
