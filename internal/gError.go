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
