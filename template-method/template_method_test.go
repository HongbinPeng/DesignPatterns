package main

import "testing"

// recordingProcessor 是一个测试用的处理器，记录流程执行顺序。
type recordingProcessor struct {
	steps []string
}

func (r *recordingProcessor) Read() string {
	r.steps = append(r.steps, "Read")
	return "data"
}

func (r *recordingProcessor) Transform(data string) string {
	r.steps = append(r.steps, "Transform")
	return data
}

func (r *recordingProcessor) Write(result string) {
	r.steps = append(r.steps, "Write")
}

// TestRunPipelineOrder 验证模板方法固定了「读->转->写」的执行顺序。
func TestRunPipelineOrder(t *testing.T) {
	r := &recordingProcessor{}
	RunPipeline(r)

	want := []string{"Read", "Transform", "Write"}
	if len(r.steps) != len(want) {
		t.Fatalf("期望执行 %d 步，实际 %d 步：%v", len(want), len(r.steps), r.steps)
	}
	for i := range want {
		if r.steps[i] != want[i] {
			t.Fatalf("第 %d 步期望 %s，实际 %s", i, want[i], r.steps[i])
		}
	}
}
