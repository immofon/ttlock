package ttlock

import "fmt"

// ErrorCode represents a TTLock API error code
type ErrorCode int

// Common errors
const (
	ErrOperationFailed       ErrorCode = 1     // 操作失败
	ErrClientIDNotExist      ErrorCode = 10000 // client_id不存在
	ErrInvalidClient         ErrorCode = 10001 // 无效的client，client_id或client_secret错
	ErrTokenNotExist         ErrorCode = 10003 // token不存在
	ErrTokenUnauthorized     ErrorCode = 10004 // token无授权，token已失效或被撤销授权
	ErrInvalidUsernameOrPass ErrorCode = 10007 // username或password错误
	ErrInvalidRefreshToken   ErrorCode = 10011 // refresh_token无效
	ErrNotLockAdmin          ErrorCode = 20002 // 不是锁管理员
	ErrInvalidUsernameFormat ErrorCode = 30002 // 用户名只能包含数字和字母
	ErrUserAlreadyExists     ErrorCode = 30003 // 用户已存在
	ErrInvalidDeleteUserID   ErrorCode = 30004 // 要删除的用户的userid不合法，只能删除当前应用注册的账号
	ErrPasswordMustBeMD5     ErrorCode = 30005 // 密码必需MD5加密
	ErrRateLimitExceeded     ErrorCode = 30006 // 超过接口调用次数限制
	ErrInvalidRequestTime    ErrorCode = 80000 // 请求时间必需为当前时间前后五分钟以内
	ErrInvalidJSONFormat     ErrorCode = 80002 // JSON格式不正确
	ErrSystemInternalError   ErrorCode = 90000 // 系统内部错误
	ErrInvalidParameter      ErrorCode = -3    // 参数不合法
	ErrPermissionDenied      ErrorCode = -2018 // 没有权限，很多接口只允许锁的最高管理员或授权管理员请求，有些接口只要有效的普通钥匙用户即可，请使用合法用户的账号和密码获取的访问令牌请求接口。
	ErrDeleteOrTransferLocks ErrorCode = -4063 // 请先删除或转移账号里所有的锁
)

// Lock related errors
const (
	ErrLockNotExist              ErrorCode = -1003 // 锁不存在
	ErrLockFrozen                ErrorCode = -2025 // 锁已被冻结，目前无法操作
	ErrCannotTransferLockToSelf  ErrorCode = -3011 // 不能把锁转移给自己
	ErrLockOperationNotSupported ErrorCode = -4043 // 此锁不支持该操作
	ErrStorageFull               ErrorCode = -4056 // 存储空间已满,操作失败
	ErrNBDeviceNotRegistered     ErrorCode = -4067 // NB设备未注册，无法发起NB操作
	ErrAutoLockTimeLimitExceeded ErrorCode = -4082 // 自动闭锁时间超限
)

// eKey related errors
const (
	ErrKeyNotExist                  ErrorCode = -1008 // 钥匙不存在
	ErrGroupNameExists              ErrorCode = -1016 // 组名已存在，请重新输入
	ErrGroupNotExist                ErrorCode = -1018 // 分组不存在
	ErrAccountBoundCannotReceiveKey ErrorCode = -1027 // 此账号已被绑定在其它账号上，无法接收电子钥匙
	ErrCannotSendKeyToSelf          ErrorCode = -2019 // 不能给自己的账号发送钥匙
	ErrCannotSendKeyToAdmin         ErrorCode = -2020 // 不能发送钥匙给管理员
	ErrCannotModifyKeyValidity      ErrorCode = -2023 // 当前不允许修改钥匙期限
	ErrReceiverNotRegistered        ErrorCode = -4064 // 发送失败，接收者账号未注册，请注册后再试
)

// Passcode related errors
const (
	ErrLockNoPasscodeData         ErrorCode = -1007 // 该锁不存在键盘密码数据
	ErrPasscodeNotExist           ErrorCode = -2009 // 密码不存在
	ErrInvalidPasscodeLength      ErrorCode = -3006 // 密码长度非法，必需为4-9位
	ErrPasscodeAlreadyExists      ErrorCode = -3007 // 已存在相同的密码，请更换
	ErrCannotModifyUnusedPasscode ErrorCode = -3008 // 无法修改从未在锁上使用过的密码
	ErrCustomPasscodeSpaceFull    ErrorCode = -3009 // 自定义密码空间已满，请删除无用密码后再试
)

// Gateway & WiFi Lock related errors
const (
	ErrNoAvailableGateway          ErrorCode = -2012 // 锁附近没有可用的网关
	ErrGatewayOffline              ErrorCode = -3002 // 网关离线，请检查后再试
	ErrGatewayBusy                 ErrorCode = -3003 // 网关正忙，请稍后再试
	ErrCannotTransferGatewayToSelf ErrorCode = -3016 // 不能把网关转移给自己
	ErrWifiLockNotConfigured       ErrorCode = -3034 // Wifi锁未配置网络，请配置网络后重试
	ErrWifiInPowerSavingMode       ErrorCode = -3035 // Wifi处于省电模式，请关闭省电模式后重试
	ErrLockOffline                 ErrorCode = -3036 // 锁离线，请检查后再试
	ErrLockBusy                    ErrorCode = -3037 // 锁正忙，请稍后再试
	ErrGatewayNotExist             ErrorCode = -4037 // 网关不存在
)

// IC Card & Fingerprint related errors
const (
	ErrICCardNotExist      ErrorCode = -1021 // 该IC卡已不存在
	ErrFingerprintNotExist ErrorCode = -1023 // 该指纹已不存在
)

var errorMessages = map[ErrorCode]string{
	// Common
	ErrOperationFailed:       "操作失败",
	ErrClientIDNotExist:      "client_id不存在",
	ErrInvalidClient:         "无效的client，client_id或client_secret错",
	ErrTokenNotExist:         "token不存在",
	ErrTokenUnauthorized:     "token无授权，token已失效或被撤销授权",
	ErrInvalidUsernameOrPass: "username或password错误",
	ErrInvalidRefreshToken:   "refresh_token无效",
	ErrNotLockAdmin:          "不是锁管理员",
	ErrInvalidUsernameFormat: "用户名只能包含数字和字母",
	ErrUserAlreadyExists:     "用户已存在",
	ErrInvalidDeleteUserID:   "要删除的用户的userid不合法，只能删除当前应用注册的账号",
	ErrPasswordMustBeMD5:     "密码必需MD5加密",
	ErrRateLimitExceeded:     "超过接口调用次数限制",
	ErrInvalidRequestTime:    "请求时间必需为当前时间前后五分钟以内",
	ErrInvalidJSONFormat:     "JSON格式不正确",
	ErrSystemInternalError:   "系统内部错误",
	ErrInvalidParameter:      "参数不合法",
	ErrPermissionDenied:      "没有权限，很多接口只允许锁的最高管理员或授权管理员请求，有些接口只要有效的普通钥匙用户即可，请使用合法用户的账号和密码获取的访问令牌请求接口。",
	ErrDeleteOrTransferLocks: "请先删除或转移账号里所有的锁",

	// Lock
	ErrLockNotExist:              "锁不存在",
	ErrLockFrozen:                "锁已被冻结，目前无法操作",
	ErrCannotTransferLockToSelf:  "不能把锁转移给自己",
	ErrLockOperationNotSupported: "此锁不支持该操作",
	ErrStorageFull:               "存储空间已满,操作失败",
	ErrNBDeviceNotRegistered:     "NB设备未注册，无法发起NB操作",
	ErrAutoLockTimeLimitExceeded: "自动闭锁时间超限",

	// eKey
	ErrKeyNotExist:                  "钥匙不存在",
	ErrGroupNameExists:              "组名已存在，请重新输入",
	ErrGroupNotExist:                "分组不存在",
	ErrAccountBoundCannotReceiveKey: "此账号已被绑定在其它账号上，无法接收电子钥匙",
	ErrCannotSendKeyToSelf:          "不能给自己的账号发送钥匙",
	ErrCannotSendKeyToAdmin:         "不能发送钥匙给管理员",
	ErrCannotModifyKeyValidity:      "当前不允许修改钥匙期限",
	ErrReceiverNotRegistered:        "发送失败，接收者账号未注册，请注册后再试",

	// Passcode
	ErrLockNoPasscodeData:         "该锁不存在键盘密码数据",
	ErrPasscodeNotExist:           "密码不存在",
	ErrInvalidPasscodeLength:      "密码长度非法，必需为4-9位",
	ErrPasscodeAlreadyExists:      "已存在相同的密码，请更换",
	ErrCannotModifyUnusedPasscode: "无法修改从未在锁上使用过的密码",
	ErrCustomPasscodeSpaceFull:    "自定义密码空间已满，请删除无用密码后再试",

	// Gateway & WiFi Lock
	ErrNoAvailableGateway:          "锁附近没有可用的网关",
	ErrGatewayOffline:              "网关离线，请检查后再试",
	ErrGatewayBusy:                 "网关正忙，请稍后再试",
	ErrCannotTransferGatewayToSelf: "不能把网关转移给自己",
	ErrWifiLockNotConfigured:       "Wifi锁未配置网络，请配置网络后重试",
	ErrWifiInPowerSavingMode:       "Wifi处于省电模式，请关闭省电模式后重试",
	ErrLockOffline:                 "锁离线，请检查后再试",
	ErrLockBusy:                    "锁正忙，请稍后再试",
	ErrGatewayNotExist:             "网关不存在",

	// IC Card & Fingerprint
	ErrICCardNotExist:      "该IC卡已不存在",
	ErrFingerprintNotExist: "该指纹已不存在",
}

// Error represents a TTLock API error
type Error struct {
	Code    ErrorCode
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("ttlock error %d: %s", e.Code, e.Message)
}

// NewError creates a new Error from a code
func NewError(code ErrorCode) *Error {
	msg, ok := errorMessages[code]
	if !ok {
		msg = "unknown error"
	}
	return &Error{
		Code:    code,
		Message: msg,
	}
}

// IsErrorCode checks if an error corresponds to a specific ErrorCode
func IsErrorCode(err error, code ErrorCode) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == code
	}
	return false
}
