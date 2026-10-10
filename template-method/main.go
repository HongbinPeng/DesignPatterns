package main

import "fmt"

// Processor 数据处理接口：定义了流程中的三个可变步骤。
// 不同的数据格式（CSV、JSON）分别实现它。
type Processor interface {
	Read() string                    // 步骤1：读取数据
	Transform(data string) string    // 步骤2：转换数据
	Write(result string)             // 步骤3：写入结果
}

// RunPipeline 模板方法：定义了固定的处理流程骨架。
//
// 注意：这里的「读取 -> 转换 -> 写入」顺序是写死的，
// 无论传进来的是 CSV 还是 JSON，都按这个顺序执行。
// 但每一步具体怎么做，由传入的 Processor 实现决定。
func RunPipeline(p Processor) {
	data := p.Read()
	result := p.Transform(data)
	p.Write(result)
}

// ---- 具体实现一：CSV 处理器 ----

type CSVProcessor struct{}

func (CSVProcessor) Read() string {
	return "name,age\n小明,18"
}

func (CSVProcessor) Transform(data string) string {
	return "【CSV转换】" + data
}

func (CSVProcessor) Write(result string) {
	fmt.Println("写入 CSV 结果：", result)
}

// ---- 具体实现二：JSON 处理器 ----

type JSONProcessor struct{}

func (JSONProcessor) Read() string {
	return `{"name":"小明","age":18}`
}

func (JSONProcessor) Transform(data string) string {
	return "【JSON转换】" + data
}

func (JSONProcessor) Write(result string) {
	fmt.Println("写入 JSON 结果：", result)
}

func main() {
	fmt.Println("== 处理 CSV ==")
	RunPipeline(CSVProcessor{})

	fmt.Println("\n== 处理 JSON ==")
	RunPipeline(JSONProcessor{})

	// 可以观察到：流程骨架（读->转->写）完全一样，
	// 但每个步骤的实际内容因格式不同而不同。
}
