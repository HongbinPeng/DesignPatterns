package main

import "fmt"

// DataSource 数据源接口：核心实现和所有装饰器都实现它。
type DataSource interface {
	Read(key string) (string, error)
}

// ---- 核心实现：内存数据源（最简单的实现） ----

// MemorySource 直接从内存 map 里读数据。
type MemorySource struct {
	data map[string]string
}

func NewMemorySource() *MemorySource {
	return &MemorySource{
		data: map[string]string{
			"name":  "小明",
			"level": "vip",
		},
	}
}

func (m *MemorySource) Read(key string) (string, error) {
	v, ok := m.data[key]
	if !ok {
		return "", fmt.Errorf("key %q 不存在", key)
	}
	return v, nil
}

// ---- 装饰器一：日志装饰器 ----

// LoggingSource 每次读取前打印一条日志，然后转发给被包装的对象。
type LoggingSource struct {
	inner DataSource // 被包装的数据源
}

func (l *LoggingSource) Read(key string) (string, error) {
	fmt.Printf("[日志] 正在读取 key：%s\n", key)
	return l.inner.Read(key) // 转发
}

// ---- 装饰器二：缓存装饰器 ----

// CachingSource 把读过的结果缓存起来，下次直接命中，不再访问 inner。
type CachingSource struct {
	inner DataSource
	cache map[string]string
}

func (c *CachingSource) Read(key string) (string, error) {
	if v, ok := c.cache[key]; ok {
		fmt.Printf("[缓存] 命中 key：%s\n", key)
		return v, nil
	}
	v, err := c.inner.Read(key)
	if err == nil {
		c.cache[key] = v
	}
	return v, err
}

func main() {
	// 核心对象
	mem := NewMemorySource()

	// 像包洋葱一样，一层层包起来：
	// 最外层是「缓存」，里面是「日志」，最里面是「内存数据源」
	source := &CachingSource{
		inner: &LoggingSource{inner: mem},
		cache: map[string]string{},
	}

	// 第一次读 name：走 缓存(miss) -> 日志 -> 内存
	if v, err := source.Read("name"); err != nil {
		fmt.Println("读取失败：", err)
	} else {
		fmt.Println("结果：", v)
	}
	fmt.Println()

	// 第二次读 name：走 缓存(hit)，不再访问底层
	if v, err := source.Read("name"); err != nil {
		fmt.Println("读取失败：", err)
	} else {
		fmt.Println("结果：", v)
	}
}
