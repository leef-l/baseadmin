package inpututil

import "github.com/gogf/gf/v2/errors/gerror"

func Require(value any) error {
	if value == nil {
		return gerror.New("请求参数不能为空")
	}
	return nil
}
