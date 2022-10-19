package protocol

// 应答结果
type Result byte

const (
	// 成功/确认
	T808_0x8001ResultSuccess = 0
	// 失败
	T808_0x8001ResultFail = 1
	// 消息有误
	T808_0x8001ResultBad = 2
	// 不支持
	T808_0x8001ResultUnsupported = 3
	// 报警处理确认
	T808_0x8001ResultAlarmConfirm = 4
)
