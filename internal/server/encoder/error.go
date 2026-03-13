package encoder

import (
	"github.com/go-kratos/kratos/v2/errors"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"net/http"
	"strings"
)

const (
	baseContentType = "application"
)

// DefaultHttpErrorResponse 默认http错误响应体
type DefaultHttpErrorResponse struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

// MakeContentType returns the content-type with base prefix.
func MakeContentType(subtype string) string {
	return strings.Join([]string{baseContentType, subtype}, "/")
}

// DefaultHttpServerErrorEncoder 默认http服务端错误响应编码
func DefaultHttpServerErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	// 拿到error并转换成 kratos Error实体
	se := errors.FromError(err)

	res := &DefaultHttpErrorResponse{
		Code:    se.Code,
		Message: se.Message,
		Success: true,
		Data:    nil,
	}

	if se.Code != http.StatusOK {
		res.Success = false
	}

	// 通过Request Header的Accept中提取出对应的编码器
	codec, _ := khttp.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", MakeContentType(codec.Name()))
	// 设置HTTP Status Code
	w.WriteHeader(int(se.Code))
	w.Write(body)
}
