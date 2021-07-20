package internal

import "net/http"

// GlobalError 全局异常结构体
type GlobalError struct {
	Status  int    `json:"-"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//获取err的内容
func (err GlobalError) Error() string {
	return err.Message
}

const (
	NoKey           = 1001 //礼品码不存在
	KeyExpired      = 1002 //礼品码过期
	UserHasEeceived = 1003 //不可重复领取
	NoGift          = 1004 //礼品全部领完
	IsEmpty         = 1005 //参数不能为空
	ValidPeriod     = 1006 //有效期不能小于当前时间
	LenFalse        = 1007 //礼品码不合法
	InternalService = 1008 //内部服务错误
)

// NoKeyError  礼品码不存在
func NoKeyError(message string) GlobalError {
	return GlobalError{
		Status:  http.StatusForbidden,
		Code:    NoKey,
		Message: message,
	}
}

// KeyExpiredError 礼品码过期
func KeyExpiredError(message string) GlobalError {
	return GlobalError{
		Status:  http.StatusForbidden,
		Code:    KeyExpired,
		Message: message,
	}
}

// UserHasEeceivedError 不可重复领取
func UserHasEeceivedError(message string) GlobalError {
	return GlobalError{
		Status:  http.StatusForbidden,
		Code:    UserHasEeceived,
		Message: message,
	}
}

// NoGiftError  礼品全部领完
func NoGiftError(message string) GlobalError {
	return GlobalError{
		Status:  http.StatusForbidden,
		Code:    NoGift,
		Message: message,
	}
}

//IsEmptyError  参数不能为空
func IsEmptyError(message string) GlobalError {
	return GlobalError{
		Status:  http.StatusForbidden,
		Code:    IsEmpty,
		Message: message,
	}
}

//ValidPeriodError  有效期不能小于当前时间
func ValidPeriodError(message string) GlobalError {
	return GlobalError{
		Status:  http.StatusForbidden,
		Code:    ValidPeriod,
		Message: message,
	}
}

//LenFalseError  礼品码不合法
func LenFalseError(message string) GlobalError {
	return GlobalError{
		Status:  http.StatusForbidden,
		Code:    LenFalse,
		Message: message,
	}
}

//InternalServiceError   内部服务错误
func InternalServiceError(message string) GlobalError {
	return GlobalError{
		Status:  http.StatusForbidden,
		Code:    InternalService,
		Message: message,
	}
}
