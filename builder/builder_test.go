package main

import "testing"

// TestBuilderDefaults 验证默认值被正确填充。
func TestBuilderDefaults(t *testing.T) {
	c := NewServerConfigBuilder().Build()

	if c.Host != "127.0.0.1" || c.Port != 8080 || c.Timeout != 30 || c.MaxConn != 100 {
		t.Fatalf("默认值不正确：%+v", c)
	}
}

// TestBuilderCustom 验证链式设置能覆盖默认值。
func TestBuilderCustom(t *testing.T) {
	c := NewServerConfigBuilder().
		Host("0.0.0.0").
		Port(9090).
		Timeout(60).
		EnableDebug().
		Build()

	if c.Host != "0.0.0.0" {
		t.Fatalf("Host 未生效：%s", c.Host)
	}
	if c.Port != 9090 {
		t.Fatalf("Port 未生效：%d", c.Port)
	}
	if !c.Debug {
		t.Fatal("EnableDebug 未生效")
	}
	// 未设置的字段应保持默认值
	if c.MaxConn != 100 {
		t.Fatalf("未设置的 MaxConn 应保持默认 100，实际 %d", c.MaxConn)
	}
}
