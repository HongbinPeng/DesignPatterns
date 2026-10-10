package main

import "fmt"

// LeaveRequest 请假请求。
type LeaveRequest struct {
	Name string
	Days int
}

// Approver 审批人接口：责任链上的每一个节点都实现它。
type Approver interface {
	// SetNext 设置下一个审批人，并返回它（方便链式构建）。
	SetNext(next Approver) Approver
	// Handle 处理请假请求。
	Handle(req *LeaveRequest)
}

// BaseApprover 基础审批人：封装了「next 字段」和「往下传」的通用逻辑。
// 具体审批人通过嵌入 BaseApprover 来复用这部分代码。
type BaseApprover struct {
	name  string   // 审批人名字
	limit int      // 自己能批准的最大天数
	next  Approver // 下一个审批人
}

// SetNext 设置下一个审批人。
func (b *BaseApprover) SetNext(next Approver) Approver {
	b.next = next
	return next
}

// Handle 通用处理逻辑：
// 如果天数在自己的权限内就批准；否则传给下一个；没有下一个就拒绝。
func (b *BaseApprover) Handle(req *LeaveRequest) {
	if req.Days <= b.limit {
		fmt.Printf("  ✅ %s 批准了 %s 的 %d 天假期\n", b.name, req.Name, req.Days)
		return
	}
	if b.next != nil {
		fmt.Printf("  ⏭️  %s 无权批准 %d 天，转交上级\n", b.name, req.Days)
		b.next.Handle(req)
	} else {
		fmt.Printf("  ❌ 申请 %d 天假期，无人能批准，已拒绝\n", req.Days)
	}
}

// ---- 具体的审批人：各自嵌入 BaseApprover，并定义权限上限 ----

// Supervisor 主管：最多批 3 天。
type Supervisor struct{ BaseApprover }

func NewSupervisor() *Supervisor {
	return &Supervisor{BaseApprover{name: "主管", limit: 3}}
}

// Manager 经理：最多批 7 天。
type Manager struct{ BaseApprover }

func NewManager() *Manager {
	return &Manager{BaseApprover{name: "经理", limit: 7}}
}

// Director 总监：最多批 15 天。
type Director struct{ BaseApprover }

func NewDirector() *Director {
	return &Director{BaseApprover{name: "总监", limit: 15}}
}

func main() {
	// 构建责任链：主管 -> 经理 -> 总监
	supervisor := NewSupervisor()
	manager := NewManager()
	director := NewDirector()

	supervisor.SetNext(manager).SetNext(director)
	// 等价于：
	// supervisor.SetNext(manager)
	// manager.SetNext(director)

	// 请求从链头（主管）开始
	fmt.Println("场景一：请 2 天假")
	supervisor.Handle(&LeaveRequest{Name: "小明", Days: 2})

	fmt.Println("\n场景二：请 5 天假")
	supervisor.Handle(&LeaveRequest{Name: "小红", Days: 5})

	fmt.Println("\n场景三：请 10 天假")
	supervisor.Handle(&LeaveRequest{Name: "小刚", Days: 10})

	fmt.Println("\n场景四：请 30 天假（超出所有人权限）")
	supervisor.Handle(&LeaveRequest{Name: "小美", Days: 30})
}
