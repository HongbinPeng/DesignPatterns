package main

import (
	"fmt"
	"sync"
)

// Config 全局配置对象。
// 整个程序只需要一份，所有地方读到的都应该是同一份。
type Config struct {
	AppName string
	Port    int
}

// 这里用「包级变量」保存唯一实例，以及保证只初始化一次的 sync.Once。
var (
	instance *Config
	once     sync.Once
)

// GetConfig 返回全局唯一的 Config 实例。
//
// 这是 Go 里实现单例最推荐的写法：借助 sync.Once，
// 无论多少个 goroutine 同时第一次调用 GetConfig，
// 初始化代码都只会执行一次，而且是线程安全的。
func GetConfig() *Config {
	once.Do(func() {
		// 只有第一次调用 GetConfig 时，才会走到这里做初始化。
		// 可以想象成：这里是「读配置文件 / 建立连接」等昂贵操作。
		instance = &Config{
			AppName: "design-patterns-demo",
			Port:    8080,
		}
	})
	return instance
}

func main() {
	// 不同地方、不同时间点拿到的是同一个对象吗？
	c1 := GetConfig()
	c2 := GetConfig()

	// c1 和 c2 指向同一块内存，所以 == 比较为 true
	fmt.Printf("c1 和 c2 是同一个实例吗？%v\n", c1 == c2)

	// 改一下 c1，c2 也会变（因为它们本来就是同一个对象）
	c1.Port = 9090
	fmt.Printf("修改 c1 后，c2 的端口 = %d（说明是同一个实例）\n", c2.Port)

	fmt.Printf("最终配置：AppName=%s, Port=%d\n", instance.AppName, instance.Port)
}
