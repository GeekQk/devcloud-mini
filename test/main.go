package main

import (
	"fmt"
	"log"
	"time"
)

type Server struct {
	host    string
	port    int
	timeout time.Duration
	maxConn int
}

func (s *Server) Start() error {
	fmt.Println(s)
	return nil
}

// Option 类型是一个函数类型，它接收一个参数：*Server
type Option func(*Server)

// 定义一系列相关返回 Option 的函数：
func WithHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}

func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func WithMaxConn(maxConn int) Option {
	return func(s *Server) {
		s.maxConn = maxConn
	}
}

// Server 的构造函数接收一个 Option 类型的不定参数
func New(options ...Option) *Server {
	svr := &Server{}
	// 遍历 Option 类型的不定参数，调用每个 Option 类型的函数
	for _, option := range options {
		option(svr)
	}
	return svr
}

func main() {
	svr := New(
		WithHost("localhost"),
		WithPort(8080),
		WithTimeout(time.Minute),
		WithMaxConn(120),
	)
	if err := svr.Start(); err != nil {
		log.Fatal(err)
	}
}
