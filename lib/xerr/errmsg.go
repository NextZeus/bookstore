package xerr

var message map[int]string

func init() {
	message = make(map[int]string)
	message[OK] = "SUCCESS"
	message[BAD_REUQEST_ERROR] = "服务器繁忙,请稍后重试"
	message[REUQES_PARAM_ERROR] = "参数错误"
	message[USER_NOT_FOUND] = "用户不存在"
}

func MapErrMsg(errCode int) string {
	if msg, ok := message[errCode]; ok {
		return msg
	} else {
		return message[BAD_REUQEST_ERROR]
	}
}