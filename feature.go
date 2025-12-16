package ttlock

import (
	"math/big"
)

// LockFeature represents a specific capability of a lock
type LockFeature int

const (
	LockFeaturePasscode           LockFeature = 0  // 支持密码
	LockFeatureICCard             LockFeature = 1  // 支持IC卡
	LockFeatureFingerprint        LockFeature = 2  // 支持指纹
	LockFeatureWristband          LockFeature = 3  // 支持Bong手环
	LockFeatureAutoLock           LockFeature = 4  // 支持自动闭锁设置
	LockFeaturePasscodeDelete     LockFeature = 5  // 密码带删除功能
	LockFeatureFirmwareUpgrade    LockFeature = 6  // 支持固件升级设置指令
	LockFeaturePasscodeManagement LockFeature = 7  // 支持密码管理功能
	LockFeatureLockCommand        LockFeature = 8  // 支持闭锁指令
	LockFeaturePasscodeVisible    LockFeature = 9  // 支持密码显示或者隐藏的控制
	LockFeatureGatewayUnlock      LockFeature = 10 // 支持网关开锁指令
	LockFeatureFreeze             LockFeature = 11 // 支持冻结、解冻锁
	LockFeatureCyclicPasscode     LockFeature = 12 // 支持循环密码
	LockFeatureDoorSensor         LockFeature = 13 // 支持门磁
	LockFeatureRemoteUnlockConfig LockFeature = 14 // 是否支持配置远程开锁功能
	LockFeatureAudioManagement    LockFeature = 15 // 支持启用或者禁用语音提示管理
	LockFeatureNB                 LockFeature = 16 // 支持NB
	// Bit 17 is deprecated (废弃（这位定义出错，废弃掉）)
	LockFeatureAdminPasscode      LockFeature = 18 // 支持读取管理员密码
	LockFeatureHotelCard          LockFeature = 19 // 支持酒店锁卡系统
	LockFeatureNoClockChip        LockFeature = 20 // 锁没有时钟芯片
	LockFeatureNoBroadcast        LockFeature = 21 // 蓝牙不广播，不能实现App点击开锁
	LockFeaturePassageMode        LockFeature = 22 // 支持某一天几点到几点常开的模式
	LockFeatureTurnOffAutoLock    LockFeature = 23 // 支持常开模式及设置自动闭锁的情况下，是否支持关闭自动闭锁
	LockFeatureWirelessKeypad     LockFeature = 24 // 支持无线键盘
	LockFeatureLightTime          LockFeature = 25 // 支持照明灯时间配置
	LockFeatureHotelCardBlacklist LockFeature = 26 // 支持酒店卡黑名单
	LockFeatureIdentityCard       LockFeature = 27 // 支持身份证
	LockFeatureTamperAlert        LockFeature = 28 // 支持防撬开关配置（启用/禁用）
	LockFeatureResetButton        LockFeature = 29 // 支持重置键配置（启用/禁用）
	LockFeaturePrivacyLock        LockFeature = 30 // 支持反锁功能配置（启用/禁用）
	// Bit 31 reserved (保留)
	LockFeatureDeadLock                LockFeature = 32  // 支持DeadLock闭锁（客户定制的闭锁模式，内外都锁住）
	LockFeaturePassageModeException    LockFeature = 33  // 支持常开模式例外
	LockFeatureCyclicICOrFingerprint   LockFeature = 34  // 支持循环指纹/卡
	LockFeaturePrivacyMode             LockFeature = 35  // App control privacy mode (indoor deadlock) (支持App控制隐私模式（内反锁）)
	LockFeatureLeftRightOpen           LockFeature = 36  // 支持左右开门设置
	LockFeatureFingerVein              LockFeature = 37  // 支持指静脉
	LockFeatureTelinkBluetooth         LockFeature = 38  // 泰凌葳蓝牙芯片
	LockFeatureNBActivation            LockFeature = 39  // 支持NB激活配置
	LockFeatureRecoverCyclicPasscode   LockFeature = 40  // 支持循环密码恢复功能
	LockFeatureWirelessKey             LockFeature = 41  // 支持无线钥匙（遥控）
	LockFeatureAccessoryBattery        LockFeature = 42  // 支持读取配件电量信息
	LockFeatureSoundVolumeLanguage     LockFeature = 43  // 支持音量及语言设置
	LockFeatureQRCode                  LockFeature = 44  // 支持二维码
	LockFeatureDoorSensorState         LockFeature = 45  // 支持门磁状态（以前也支持，但没有门磁未知状态，增加未知状态，使得门磁状态更加准确）
	LockFeaturePassageModeAutoUnlock   LockFeature = 46  // 支持常开模式自动开锁设置
	LockFeatureFingerprintDistribution LockFeature = 47  // Must be 1 if fingerprint distribution is supported (支持指纹下发功能（为了Web页面简化显示，支持指纹下发的，不管是哪种指纹，这一位都必须要设置为1）)
	LockFeatureZhongZhengFingerprint   LockFeature = 48  // 支持中正指纹下发功能
	LockFeatureShengYuanFingerprint    LockFeature = 49  // 支持晟元指纹下发功能
	LockFeatureWirelessDoorSensor      LockFeature = 50  // 支持无线门磁
	LockFeatureDoorUnclosedAlarm       LockFeature = 51  // 支持门未关报警
	LockFeatureProximitySensor         LockFeature = 52  // 支持接近感应
	LockFeature3DFace                  LockFeature = 53  // 支持3D人脸
	LockFeatureAutoLockPairing         LockFeature = 54  // 支持全自动锁配
	LockFeatureCPUCard                 LockFeature = 55  // 支持CPU卡
	LockFeatureWiFi                    LockFeature = 56  // 支持WiFi
	LockFeatureWiFiStaticIP            LockFeature = 58  // WiFi锁支持固定IP地址
	LockFeatureIncompletePasscode      LockFeature = 60  // 支持不完整密码锁
	LockFeatureDoubleAuth              LockFeature = 63  // 支持双重认证
	LockFeatureXiongMaiVideo           LockFeature = 67  // 支持雄迈可视对讲功能
	LockFeatureZhiAnFace               LockFeature = 69  // 支持指安人脸下发
	LockFeaturePalmVein                LockFeature = 70  // 支持掌静脉
	LockFeatureOneTimeQRCode           LockFeature = 74  // 支持单次二维码
	LockFeatureThirdPartyBluetooth     LockFeature = 77  // 支持第三方蓝牙设备接入
	LockFeatureWiFiPowerSave           LockFeature = 83  // 支持WiFi省电时间段配置
	LockFeatureMultifuncWirelessKeypad LockFeature = 84  // 支持多功能无线键盘
	LockFeatureCustomQRCode            LockFeature = 108 // 支持自定义二维码
)

func (f LockFeature) String() string {
	switch f {
	case LockFeaturePasscode:
		return "支持密码"
	case LockFeatureICCard:
		return "支持IC卡"
	case LockFeatureFingerprint:
		return "支持指纹"
	case LockFeatureWristband:
		return "支持Bong手环"
	case LockFeatureAutoLock:
		return "支持自动闭锁设置"
	case LockFeaturePasscodeDelete:
		return "密码带删除功能"
	case LockFeatureFirmwareUpgrade:
		return "支持固件升级设置指令"
	case LockFeaturePasscodeManagement:
		return "支持密码管理功能"
	case LockFeatureLockCommand:
		return "支持闭锁指令"
	case LockFeaturePasscodeVisible:
		return "支持密码显示或者隐藏的控制"
	case LockFeatureGatewayUnlock:
		return "支持网关开锁指令"
	case LockFeatureFreeze:
		return "支持冻结、解冻锁"
	case LockFeatureCyclicPasscode:
		return "支持循环密码"
	case LockFeatureDoorSensor:
		return "支持门磁"
	case LockFeatureRemoteUnlockConfig:
		return "是否支持配置远程开锁功能"
	case LockFeatureAudioManagement:
		return "支持启用或者禁用语音提示管理"
	case LockFeatureNB:
		return "支持NB"
	case LockFeatureAdminPasscode:
		return "支持读取管理员密码"
	case LockFeatureHotelCard:
		return "支持酒店锁卡系统"
	case LockFeatureNoClockChip:
		return "锁没有时钟芯片"
	case LockFeatureNoBroadcast:
		return "蓝牙不广播，不能实现App点击开锁"
	case LockFeaturePassageMode:
		return "支持某一天几点到几点常开的模式"
	case LockFeatureTurnOffAutoLock:
		return "支持常开模式及设置自动闭锁的情况下，是否支持关闭自动闭锁"
	case LockFeatureWirelessKeypad:
		return "支持无线键盘"
	case LockFeatureLightTime:
		return "支持照明灯时间配置"
	case LockFeatureHotelCardBlacklist:
		return "支持酒店卡黑名单"
	case LockFeatureIdentityCard:
		return "支持身份证"
	case LockFeatureTamperAlert:
		return "支持防撬开关配置（启用/禁用）"
	case LockFeatureResetButton:
		return "支持重置键配置（启用/禁用）"
	case LockFeaturePrivacyLock:
		return "支持反锁功能配置（启用/禁用）"
	case LockFeatureDeadLock:
		return "支持DeadLock闭锁（客户定制的闭锁模式，内外都锁住）"
	case LockFeaturePassageModeException:
		return "支持常开模式例外"
	case LockFeatureCyclicICOrFingerprint:
		return "支持循环指纹/卡"
	case LockFeaturePrivacyMode:
		return "支持App控制隐私模式（内反锁）"
	case LockFeatureLeftRightOpen:
		return "支持左右开门设置"
	case LockFeatureFingerVein:
		return "支持指静脉"
	case LockFeatureTelinkBluetooth:
		return "泰凌葳蓝牙芯片"
	case LockFeatureNBActivation:
		return "支持NB激活配置"
	case LockFeatureRecoverCyclicPasscode:
		return "支持循环密码恢复功能"
	case LockFeatureWirelessKey:
		return "支持无线钥匙（遥控）"
	case LockFeatureAccessoryBattery:
		return "支持读取配件电量信息"
	case LockFeatureSoundVolumeLanguage:
		return "支持音量及语言设置"
	case LockFeatureQRCode:
		return "支持二维码"
	case LockFeatureDoorSensorState:
		return "支持门磁状态（增加未知状态，使得门磁状态更加准确）"
	case LockFeaturePassageModeAutoUnlock:
		return "支持常开模式自动开锁设置"
	case LockFeatureFingerprintDistribution:
		return "支持指纹下发功能（Web显示简化，支持指纹下发的这一位必须设置为1）"
	case LockFeatureZhongZhengFingerprint:
		return "支持中正指纹下发功能"
	case LockFeatureShengYuanFingerprint:
		return "支持晟元指纹下发功能"
	case LockFeatureWirelessDoorSensor:
		return "支持无线门磁"
	case LockFeatureDoorUnclosedAlarm:
		return "支持门未关报警"
	case LockFeatureProximitySensor:
		return "支持接近感应"
	case LockFeature3DFace:
		return "支持3D人脸"
	case LockFeatureAutoLockPairing:
		return "支持全自动锁配"
	case LockFeatureCPUCard:
		return "支持CPU卡"
	case LockFeatureWiFi:
		return "支持WiFi"
	case LockFeatureWiFiStaticIP:
		return "WiFi锁支持固定IP地址"
	case LockFeatureIncompletePasscode:
		return "支持不完整密码锁"
	case LockFeatureDoubleAuth:
		return "支持双重认证"
	case LockFeatureXiongMaiVideo:
		return "支持雄迈可视对讲功能"
	case LockFeatureZhiAnFace:
		return "支持指安人脸下发"
	case LockFeaturePalmVein:
		return "支持掌静脉"
	case LockFeatureOneTimeQRCode:
		return "支持单次二维码"
	case LockFeatureThirdPartyBluetooth:
		return "支持第三方蓝牙设备接入"
	case LockFeatureWiFiPowerSave:
		return "支持WiFi省电时间段配置"
	case LockFeatureMultifuncWirelessKeypad:
		return "支持多功能无线键盘"
	case LockFeatureCustomQRCode:
		return "支持自定义二维码"
	default:
		return "未知功能"
	}
}

// HasFeature checks if the given feature value string supports the specified feature.
// The featureValue is a hexadecimal string.
func HasFeature(featureValue string, feature LockFeature) bool {
	if featureValue == "" {
		return false
	}
	val := new(big.Int)
	// Parse the hex string
	val.SetString(featureValue, 16)
	// Check the bit at the feature position
	return val.Bit(int(feature)) == 1
}

// SupportsFeature checks if the lock supports the specified feature.
func (l *Lock) SupportsFeature(feature LockFeature) bool {
	return HasFeature(l.FeatureValue, feature)
}

// SupportsFeature checks if the lock detail supports the specified feature.
func (l *LockDetail) SupportsFeature(feature LockFeature) bool {
	return HasFeature(l.FeatureValue, feature)
}
