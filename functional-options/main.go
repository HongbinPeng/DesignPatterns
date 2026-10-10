package main

import "fmt"

// Option 是一个函数类型：接收 *Server，并修改它。
// 这是整个模式的核心——把「修改动作」当成一个可以传参的值。
type Option func(*Server)

// Server 服务配置。
type Server struct {
	Host    string
	Port    int
	Timeout int
	MaxConn int
}

// ---- 一组 WithXXX 函数，每个返回一个 Option ----

// WithHost 返回一个「把 Host 设为指定值」的 Option。
func WithHost(host string) Option {
	return func(s *Server) {
		s.Host = host
	}
}

func WithPort(port int) Option {
	return func(s *Server) {
		s.Port = port
	}
}

func WithTimeout(seconds int) Option {
	return func(s *Server) {
		s.Timeout = seconds
	}
}

func WithMaxConn(n int) Option {
	return func(s *Server) {
		s.MaxConn = n
	}
}

// NewServer 构造函数：接受任意数量的 Option（可变参数）。
//
// 先填好默认值，再依次执行传入的 option 去覆盖默认值。
func NewServer(opts ...Option) *Server {
	s := &Server{
		Host:    "127.0.0.1",
		Port:    8080,
		Timeout: 30,
		MaxConn: 100,
	}
	for _, opt := range opts {
		opt(s) // 应用每一个修改动作
	}
	return s
}

func main() {
	// 场景一：不传任何 option，全用默认值
	s1 := NewServer()
	fmt.Printf("默认：%+v\n", s1)

	// 场景二：只改一个
	s2 := NewServer(WithPort(9090))
	fmt.Printf("只改端口：%+v\n", s2)

	// 场景三：组合多个，顺序随意
	s3 := NewServer(WithHost("0.0.0.0"), WithTimeout(60), WithMaxConn(1000))
	fmt.Printf("组合多个：%+v\n", s3)
}
