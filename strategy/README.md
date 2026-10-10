# 策略模式（Strategy）

> **一句话**：把「一段会变化的算法/规则」抽出来，做成可随时替换的「策略」，让调用方灵活切换。

## 生活中的类比 🚗

你开车去公司，导航给你三条路线：**最快路线、躲避拥堵、少收费**。

目的地（目标）不变，但「怎么走」（策略）可以根据当时情况随时切换。而且切换策略这件事，对「开车去公司」这个大目标来说是透明的——你只是换了个走法，车还是那辆车，目的地还是那个目的地。

**策略模式，就是「同一件事，多种干法，随时可换」。**

## 它解决什么问题？

看这段让人头疼的「if-else 全家桶」：

```go
func calcPrice(userType string, price float64) float64 {
    if userType == "normal" {
        return price
    } else if userType == "member" {
        return price * 0.9
    } else if userType == "vip" {
        return price * 0.8
    }
    return price
}
```

问题：
1. **规则越加越多**，if-else 越来越长，越来越难读
2. **加一种新用户**，就要改这个函数（违反开闭原则）
3. **规则散落**，同一个折扣逻辑可能被复制粘贴到多处

策略模式把「每种规则」变成一个独立的策略对象，调用方（上下文）持有一个策略，随时可以换。

## 应用场景（Go 后端）

- **计价/折扣**：不同会员等级、不同活动（满减/打折/秒杀）
- **支付路由**：按渠道、金额、风控动态选支付方式
- **灰度/分流**：按用户、地区、版本走不同逻辑
- **排序、压缩、序列化算法**：同一接口，多种实现按需切换

## 优势 ✅

1. **消除 if-else 地狱**：每种策略独立成一个类型，各管各的
2. **符合开闭原则**：加新策略 = 加个新类型，不改老代码
3. **可动态切换**：运行时随时 `SetStrategy`，灵活应对变化
4. **易测试**：每个策略可以单独测试

## Go 里怎么写？

三步：

```go
// 1. 定义策略接口
type DiscountStrategy interface {
    Discount(price float64) float64
}

// 2. 每种策略一个实现
type MemberStrategy struct{}
func (MemberStrategy) Discount(price float64) float64 { return price * 0.9 }

// 3. 上下文持有策略，并允许切换
type PriceCalculator struct { strategy DiscountStrategy }
func (c *PriceCalculator) SetStrategy(s DiscountStrategy) { c.strategy = s }
func (c *PriceCalculator) Calculate(price float64) float64 { return c.strategy.Discount(price) }
```

## 策略模式 vs 工厂模式（高频混淆点！）

这俩经常一起出现，但职责不同，一定要分清：

| | 工厂模式 | 策略模式 |
| --- | --- | --- |
| 核心动作 | **创建**对象 | **替换**行为/算法 |
| 关注点 | 谁来 new、new 哪个 | 用哪种规则去计算 |
| 典型搭配 | 工厂负责「按类型造出策略」 | 策略负责「提供可替换的算法」 |

> 实战里经常是**两个一起用**：用工厂模式根据配置 `new` 出对应的策略，再用策略模式去执行。别慌，先各自理解清楚，组合是水到渠成的事。

## 何时用 / 何时不用

**用**：一段「同一接口、多种实现」的算法/规则，需要动态切换或频繁新增。

**别用**：就一两种固定分支、未来也不会多，直接 if-else 反而更直白。

## 一句话总结

> 「同一件事有多种干法」→ 把每种干法抽成策略，谁要换谁就换。
