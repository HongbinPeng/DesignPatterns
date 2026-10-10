# 建造者模式（Builder）

> **一句话**：把「构造一个**参数很多、很复杂**的对象」的过程，拆成一步步来，最后 `Build()` 出成品。

## 生活中的类比 🍔

点一个自选汉堡：面包、牛肉饼、生菜、芝士、酱料……你可以**一步一步**加，加什么、加几层都自己定，最后店员给你打包成一个完整汉堡。

建造者模式就是「自选汉堡的点餐流程」：不要求你一次性把所有配料都填对，而是让你按需、按顺序、可读地一步步搭建，最后统一产出。

## 它解决什么问题？

看看这个「噩梦构造函数」：

```go
// 参数这么多，还都是同类型，调用时谁分得清哪个是哪个？
NewServer("127.0.0.1", 8080, 30, 100, true, false, 60)
```

问题：
1. **参数一多就分不清**：`30` 是超时还是并发？`true` 是 debug 还是 tls？
2. **没法优雅地设置默认值**：想「只改超时、其他用默认」，靠构造函数几乎不可能
3. **扩展参数就是灾难**：多加一个参数，所有调用处都要跟着改

建造者模式用「**链式调用 + 语义化方法名 + 最后统一 Build**」来化解这些问题。

## 应用场景（Go 后端）

- 构造**参数超多**的配置对象（服务配置、HTTP 客户端配置）
- 拼装**复杂 SQL 查询**（各种 where/join/order by 组合）
- 构造**复杂请求**（带各种 header、query、body 的 HTTP 请求）
- 任何「可选字段很多、需要默认值」的对象

## 优势 ✅

1. **可读性极强**：`.Timeout(30).MaxConn(100)` 一眼看懂每个值是什么
2. **默认值优雅**：Builder 里先填好默认值，调用方只覆盖想改的
3. **链式流畅**：按需调用，顺序自由，代码像在「描述」对象而不是「填坑」
4. **扩展友好**：加新字段只需给 Builder 加个方法，不影响老调用

## Go 里怎么写？

```go
type ServerConfig struct {
    Host    string
    Port    int
    Timeout int
}

// Builder 里先放好默认值
func NewServerConfigBuilder() *ServerConfigBuilder {
    return &ServerConfigBuilder{config: &ServerConfig{
        Host: "127.0.0.1", Port: 8080, Timeout: 30,
    }}
}

// 每个字段一个语义化方法，返回 builder 自身以便链式调用
func (b *ServerConfigBuilder) Host(h string) *ServerConfigBuilder {
    b.config.Host = h
    return b
}

func (b *ServerConfigBuilder) Build() *ServerConfig {
    return b.config
}
```

调用时：

```go
config := NewServerConfigBuilder().
    Host("0.0.0.0").
    Port(9090).
    Build()   // 最终成品
```

## 和「函数选项模式」的区别（重要！）

这俩在 Go 里长得有点像，都用来解决「参数多」的问题，但思路不同：

- **建造者**：先 `NewBuilder()`，再一步步 `.XXX()` 设置，最后 `.Build()`。适合**步骤感强、需要复用中间状态**的场景。
- **函数选项**：`NewServer(WithHost("x"), WithPort(80))`，把可选参数直接作为构造函数的可变参数。**更简洁、更 Go**，是 Go 社区最主流的做法。

> 一个实用建议：Go 项目里，**优先考虑函数选项模式**（见下一个文件夹）。建造者模式在 Go 里相对少见，但理解它有助于你看懂其他语言（Java 等）的代码和思路。

## 何时用 / 何时不用

**用**：对象构造复杂、可选参数多、需要默认值、构造过程有步骤顺序。

**别用**：就两三个参数，直接构造就行，套 Builder 反而啰嗦。

## 一句话总结

> 参数多得吓人的构造函数 → 换成「链式设置 + Build 成品」。
