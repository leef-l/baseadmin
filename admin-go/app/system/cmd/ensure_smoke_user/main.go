package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	"gbaseadmin/app/system/internal/ops/smokeuser"
	"gbaseadmin/internal/cmdutil"
)

func main() {
	os.Exit(run())
}

func run() int {
	var (
		configPath string
		username   string
		password   string
		nickname   string
	)

	flag.StringVar(&configPath, "config", "", "GoFrame config file path")
	flag.StringVar(&username, "username", "", "smoke username")
	flag.StringVar(&password, "password", "", "smoke password")
	flag.StringVar(&nickname, "nickname", "", "smoke nickname")
	flag.Parse()

	if err := cmdutil.UseConfigFile(configPath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	result, err := smokeuser.Ensure(context.Background(), smokeuser.EnsureOptions{
		Username: username,
		Password: password,
		Nickname: nickname,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	fmt.Printf(
		"smoke user ready: role_id=%d user_id=%d role_created=%t role_updated=%t user_created=%t user_updated=%t menu_count=%d\n",
		result.RoleID,
		result.UserID,
		result.RoleCreated,
		result.RoleUpdated,
		result.UserCreated,
		result.UserUpdated,
		result.MenuCount,
	)
	return 0
}
