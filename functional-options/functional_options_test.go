package main

import "testing"

// TestNewServerDefaults 验证不传 option 时使用默认值。
func TestNewServerDefaults(t *testing.T) {
	s := NewServer()

	if s.Host != "127.0.0.1" || s.Port != 8080 || s.Timeout != 30 || s.MaxConn != 100 {
		t.Fatalf("默认值不正确：%+v", s)
	}
}

// TestNewServerWithOptions 验证 option 能覆盖默认值。
func TestNewServerWithOptions(t *testing.T) {
	s := NewServer(WithPort(9090), WithTimeout(60))

	if s.Port != 9090 {
		t.Fatalf("WithPort 未生效：%d", s.Port)
	}
	if s.Timeout != 60 {
		t.Fatalf("WithTimeout 未生效：%d", s.Timeout)
	}
	// 未覆盖的字段保持默认
	if s.Host != "127.0.0.1" {
		t.Fatalf("未设置的 Host 应保持默认，实际 %s", s.Host)
	}
}
