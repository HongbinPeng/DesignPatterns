package main

import "fmt"

// ServerConfig 服务配置，字段很多，而且都有合理的默认值。
type ServerConfig struct {
	Host     string // 监听地址
	Port     int    // 端口
	Timeout  int    // 超时时间（秒）
	MaxConn  int    // 最大连接数
	Debug    bool   // 是否开启调试
	CertFile string // TLS 证书文件（可选）
}

// ServerConfigBuilder 建造者：负责一步步搭出 ServerConfig。
type ServerConfigBuilder struct {
	config *ServerConfig
}

// NewServerConfigBuilder 返回一个 Builder，并**预先填好默认值**。
// 这样调用方只需要覆盖自己想改的那几个字段。
func NewServerConfigBuilder() *ServerConfigBuilder {
	return &ServerConfigBuilder{
		config: &ServerConfig{
			Host:    "127.0.0.1",
			Port:    8080,
			Timeout: 30,
			MaxConn: 100,
			Debug:   false,
		},
	}
}

// 下面每个字段一个语义化方法，方法名就是字段名，
// 一眼就能看出在设置什么。返回 builder 自身，实现链式调用。

func (b *ServerConfigBuilder) Host(host string) *ServerConfigBuilder {
	b.config.Host = host
	return b
}

func (b *ServerConfigBuilder) Port(port int) *ServerConfigBuilder {
	b.config.Port = port
	return b
}

func (b *ServerConfigBuilder) Timeout(seconds int) *ServerConfigBuilder {
	b.config.Timeout = seconds
	return b
}

func (b *ServerConfigBuilder) MaxConn(n int) *ServerConfigBuilder {
	b.config.MaxConn = n
	return b
}

func (b *ServerConfigBuilder) EnableDebug() *ServerConfigBuilder {
	b.config.Debug = true
	return b
}

func (b *ServerConfigBuilder) CertFile(path string) *ServerConfigBuilder {
	b.config.CertFile = path
	return b
}

// Build 返回最终搭建好的配置对象。
func (b *ServerConfigBuilder) Build() *ServerConfig {
	return b.config
}

func main() {
	// 场景一：全部用默认值
	defaultConfig := NewServerConfigBuilder().Build()
	fmt.Printf("默认配置：%+v\n\n", defaultConfig)

	// 场景二：链式设置，只改自己关心的字段
	customConfig := NewServerConfigBuilder().
		Host("0.0.0.0").
		Port(9090).
		Timeout(60).
		MaxConn(1000).
		EnableDebug().
		Build()
	fmt.Printf("自定义配置：%+v\n", customConfig)
}
