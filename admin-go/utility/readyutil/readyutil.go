package readyutil

import "github.com/gogf/gf/v2/frame/g"

func CheckDatabase() error {
	return g.DB().PingMaster()
}
