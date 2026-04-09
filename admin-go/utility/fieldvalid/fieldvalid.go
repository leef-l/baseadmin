package fieldvalid

import "github.com/gogf/gf/v2/errors/gerror"

func Enum(name string, value int, allowed ...int) error {
	for _, current := range allowed {
		if value == current {
			return nil
		}
	}
	return gerror.New(name + "值不合法")
}

func Binary(name string, value int) error {
	return Enum(name, value, 0, 1)
}

func NonNegative(name string, value int) error {
	if value < 0 {
		return gerror.New(name + "不能小于0")
	}
	return nil
}

func NonNegative64(name string, value int64) error {
	if value < 0 {
		return gerror.New(name + "不能小于0")
	}
	return nil
}
