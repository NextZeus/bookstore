package xhttp

import (
	"bookstore/lib/xerr"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
	"net/http"
)

//http方法
func HttpResult(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	if err == nil {
		// 成功返回
		r := Success(resp)
		httpx.WriteJson(w, http.StatusOK, r)
	} else {
		// 错误返回
		errcode := xerr.BAD_REUQEST_ERROR
		errmsg := xerr.MapErrMsg(errcode)

		if  e, ok := err.(*xerr.CodeError); ok {
			errcode = e.GetErrCode()
			errmsg = e.GetErrMsg()
		} else {
			originErr := errors.Cause(err) // err类型
			if gstatus, ok := status.FromError(originErr); ok {
				// grpc err错误
				errmsg = gstatus.Message()
			}
		}
		logx.WithContext(r.Context()).Error("【GATEWAY-SRV-ERR】 : %+v ",err)

		httpx.WriteJson(w, http.StatusBadRequest, Error(errcode, errmsg))
	}
}

// http参数错误返回
func ParamErrorResult(r *http.Request, w http.ResponseWriter, err error) {
	errMsg := fmt.Sprintf("%s, %s", xerr.MapErrMsg(xerr.REUQES_PARAM_ERROR), err.Error())
	httpx.WriteJson(w, http.StatusBadRequest, Error(xerr.REUQES_PARAM_ERROR, errMsg))
}