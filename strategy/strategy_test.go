package main

import "testing"

func TestStrategyDiscount(t *testing.T) {
	cases := []struct {
		name     string
		strategy DiscountStrategy
		price    float64
		want     float64
	}{
		{"普通用户", &NormalStrategy{}, 100, 100},
		{"会员", &MemberStrategy{}, 100, 90},
		{"VIP", &VIPStrategy{}, 100, 80},
	}

	for _, c := range cases {
		calc := NewPriceCalculator(c.strategy)
		if got := calc.Calculate(c.price); got != c.want {
			t.Fatalf("%s：期望 %.2f，实际 %.2f", c.name, c.want, got)
		}
	}
}

// TestStrategySwitch 验证运行时能动态切换策略。
func TestStrategySwitch(t *testing.T) {
	calc := NewPriceCalculator(&NormalStrategy{})
	if got := calc.Calculate(100); got != 100 {
		t.Fatalf("普通用户期望 100，实际 %v", got)
	}

	calc.SetStrategy(&VIPStrategy{})
	if got := calc.Calculate(100); got != 80 {
		t.Fatalf("切换 VIP 后期望 80，实际 %v", got)
	}
}
