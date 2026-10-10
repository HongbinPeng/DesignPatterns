package main

import "fmt"

// SmsSender 你项目里统一使用的短信发送接口。
// 业务代码只认这个接口，不关心背后是阿里云还是腾讯云。
type SmsSender interface {
	Send(phone, content string) error
}

// ---- 第三方 SDK（接口各不相同，而且你改不了它们） ----

// AliyunSmsSDK 阿里云短信 SDK：方法叫 SendMessage，参数叫 mobile/text。
type AliyunSmsSDK struct{}

func (a *AliyunSmsSDK) SendMessage(mobile, text string) error {
	fmt.Printf("  [阿里云 SDK] 发送短信到 %s：%s\n", mobile, text)
	return nil
}

// TencentSmsSDK 腾讯云短信 SDK：方法叫 PushSms，返回值多了一个 bool。
type TencentSmsSDK struct{}

func (t *TencentSmsSDK) PushSms(phoneNumber, body string) (bool, error) {
	fmt.Printf("  [腾讯云 SDK] 发送短信到 %s：%s\n", phoneNumber, body)
	return true, nil
}

// ---- 适配器：把第三方 SDK「翻译」成统一的 SmsSender 接口 ----

// AliyunAdapter 把阿里云 SDK 适配成 SmsSender。
type AliyunAdapter struct {
	sdk *AliyunSmsSDK
}

// Send 实现 SmsSender 接口，内部把参数转成阿里云 SDK 的调用。
func (a *AliyunAdapter) Send(phone, content string) error {
	return a.sdk.SendMessage(phone, content)
}

// TencentAdapter 把腾讯云 SDK 适配成 SmsSender。
type TencentAdapter struct {
	sdk *TencentSmsSDK
}

// Send 实现 SmsSender 接口，注意要处理腾讯云 SDK 多出来的 bool 返回值。
func (t *TencentAdapter) Send(phone, content string) error {
	_, err := t.sdk.PushSms(phone, content)
	return err
}

// ---- 业务代码：只跟 SmsSender 接口打交道 ----

// SendVerifyCode 发送验证码的业务逻辑。
// 注意：它完全不知道、也不关心底层是哪个供应商，换供应商它一行都不用改。
func SendVerifyCode(sender SmsSender, phone, code string) {
	fmt.Println("发送验证码...")
	if err := sender.Send(phone, "您的验证码是 "+code); err != nil {
		fmt.Println("发送失败：", err)
		return
	}
	fmt.Println("发送成功！")
}

func main() {
	// 用阿里云
	aliyun := &AliyunAdapter{sdk: &AliyunSmsSDK{}}
	SendVerifyCode(aliyun, "13800138000", "123456")
	fmt.Println()

	// 换腾讯云：只需换一个适配器，业务逻辑 SendVerifyCode 完全不动
	tencent := &TencentAdapter{sdk: &TencentSmsSDK{}}
	SendVerifyCode(tencent, "13900139000", "654321")
}
