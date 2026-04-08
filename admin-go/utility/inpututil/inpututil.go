package inpututil

import (
	"reflect"

	"github.com/gogf/gf/v2/errors/gerror"
)

func Require(value any) error {
	if value == nil {
		return gerror.New("请求参数不能为空")
	}
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		if rv.IsNil() {
			return gerror.New("请求参数不能为空")
		}
	}
	return nil
}
