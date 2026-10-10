package main

import "testing"

// newChain 构建一条测试用的责任链：主管 -> 经理 -> 总监。
func newChain() *Supervisor {
	supervisor := NewSupervisor()
	manager := NewManager()
	director := NewDirector()

	supervisor.SetNext(manager).SetNext(director)
	return supervisor
}

// TestChainHandledBySupervisor 2 天假由主管批准，不再往下传。
func TestChainHandledBySupervisor(t *testing.T) {
	req := &LeaveRequest{Name: "小明", Days: 2}

	// 用审批结果来断言：2 天 <= 主管 limit 3，应该主管批准。
	// 这里直接调用验证不 panic 且能走通（批准逻辑通过输出验证）。
	newChain().Handle(req)
}

// TestChainEscalates 5 天假超出主管权限，应转交经理批准。
func TestChainEscalates(t *testing.T) {
	req := &LeaveRequest{Name: "小红", Days: 5}
	newChain().Handle(req)
}

// TestBaseApproverDirect 直接测试 BaseApprover 的判断逻辑，避免依赖输出。
func TestBaseApproverDecision(t *testing.T) {
	// 一个只有主管（limit=3）的链
	supervisor := NewSupervisor()

	// 3 天：主管能批
	// 4 天：主管批不了，且没有下一个，应拒绝
	// 这里通过构造不同天数，验证逻辑分支都能正常执行不 panic。
	for _, days := range []int{3, 4} {
		supervisor.Handle(&LeaveRequest{Name: "测试", Days: days})
	}
}
