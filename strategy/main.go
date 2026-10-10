package main

import "fmt"

// DiscountStrategy 折扣策略接口：所有折扣算法都实现它。
type DiscountStrategy interface {
	Discount(price float64) float64
}

// ---- 各种折扣策略 ----

// NormalStrategy 普通用户：不打折。
type NormalStrategy struct{}

func (NormalStrategy) Discount(price float64) float64 {
	return price
}

// MemberStrategy 会员：9 折。
type MemberStrategy struct{}

func (MemberStrategy) Discount(price float64) float64 {
	return price * 0.9
}

// VIPStrategy VIP：8 折。
type VIPStrategy struct{}

func (VIPStrategy) Discount(price float64) float64 {
	return price * 0.8
}

// PriceCalculator 计价上下文：持有一个折扣策略，并允许运行时切换。
type PriceCalculator struct {
	strategy DiscountStrategy
}

func NewPriceCalculator(s DiscountStrategy) *PriceCalculator {
	return &PriceCalculator{strategy: s}
}

// SetStrategy 切换策略（策略模式的关键：行为可动态替换）。
func (c *PriceCalculator) SetStrategy(s DiscountStrategy) {
	c.strategy = s
}

// Calculate 用当前策略计算折后价。
func (c *PriceCalculator) Calculate(price float64) float64 {
	return c.strategy.Discount(price)
}

func main() {
	price := 100.0

	calculator := NewPriceCalculator(&NormalStrategy{})
	fmt.Printf("普通用户应付：%.2f\n", calculator.Calculate(price))

	// 用户升级成 VIP，只需切换策略，计价逻辑不用改
	calculator.SetStrategy(&VIPStrategy{})
	fmt.Printf("VIP 用户应付：%.2f\n", calculator.Calculate(price))

	// 再换成会员
	calculator.SetStrategy(&MemberStrategy{})
	fmt.Printf("会员应付：%.2f\n", calculator.Calculate(price))
}
